package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

const testPass = "test-password-123"

// newTestApp runs against a copy-free temp data dir seeded from the real site (SITE_DIR = repo root).
func newTestApp(t *testing.T) (*App, *httptest.Server) {
	t.Helper()
	app, err := NewApp(Config{SiteDir: "..", DataDir: t.TempDir(), AdminUser: "admin", AdminPass: testPass, Secret: []byte("s3cret")})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(app.Routes())
	t.Cleanup(srv.Close)
	return app, srv
}

func TestSeedKeepsAllSiteData(t *testing.T) {
	app, _ := newTestApp(t)
	c := app.store.Get()
	if len(c.Team) == 0 || len(c.Films) == 0 || len(c.Gallery) == 0 { // exact fidelity is checked against the seed below
		t.Fatalf("counts team=%d films=%d gallery=%d", len(c.Team), len(c.Films), len(c.Gallery))
	}
	seed, err := readDataJS("../assets/data.js")
	if err != nil {
		t.Fatal(err)
	}
	js := app.store.DataJS()
	if !bytes.HasPrefix(js, []byte("window.NS_DATA = ")) {
		t.Fatalf("data.js prefix: %q", js[:30])
	}
	var want, got any
	json.Unmarshal(seed, &want)
	json.Unmarshal(bytes.TrimSuffix(bytes.TrimPrefix(js, []byte("window.NS_DATA = ")), []byte(";\n")), &got)
	if !reflect.DeepEqual(want, got) {
		t.Fatal("served data.js differs from the seed: some field was lost in the Go model")
	}
}

func TestUpdateHistoryRestore(t *testing.T) {
	app, _ := newTestApp(t)
	before := len(app.store.Get().Team)
	if err := app.store.Update(func(c *Content) error { c.Team = c.Team[1:]; return nil }); err != nil {
		t.Fatal(err)
	}
	if got := len(app.store.Get().Team); got != before-1 {
		t.Fatalf("team after delete = %d", got)
	}
	snaps, err := app.store.History()
	if err != nil || len(snaps) != 1 {
		t.Fatalf("history = %v, %v", snaps, err)
	}
	if err := app.store.Restore(snaps[0].Name); err != nil {
		t.Fatal(err)
	}
	if got := len(app.store.Get().Team); got != before {
		t.Fatalf("team after restore = %d", got)
	}
	if err := app.store.Restore("../content.json"); err == nil {
		t.Fatal("restore accepted a path outside history")
	}
	// survives a restart
	again, err := OpenStore(app.cfg.DataDir, "unused")
	if err != nil || len(again.Get().Team) != before {
		t.Fatalf("reopen: %v team=%d", err, len(again.Get().Team))
	}
}

func TestSiteServesOnlyPublicFiles(t *testing.T) {
	_, srv := newTestApp(t)
	cases := map[string]int{
		"/": 200, "/about": 200, "/about.html": 200, "/assets/app.js": 200, "/assets/data.js": 200, "/favicon.ico": 200,
		"/server/main.go": 404, "/server/": 404, "/.git/config": 404, "/assets/": 404, "/data/content.json": 404,
		"/vercel.json": 404, "/README.md": 404, "/nope": 404, "/media/": 404, "/upload/images/../../server/main.go": 404,
	}
	for p, want := range cases {
		res, err := http.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != want {
			t.Errorf("GET %s = %d, want %d", p, res.StatusCode, want)
		}
	}
}

func TestAdminFlow(t *testing.T) {
	app, srv := newTestApp(t)
	jar, _ := cookiejar.New(nil)
	cl := &http.Client{Jar: jar}

	// not logged in -> login page
	res, _ := cl.Get(srv.URL + "/admin/team")
	if res.Request.URL.Path != "/admin/login" {
		t.Fatalf("unauthenticated landed on %s", res.Request.URL.Path)
	}
	res.Body.Close()

	// wrong password
	res, _ = cl.PostForm(srv.URL+"/admin/login", url.Values{"user": {"admin"}, "password": {"nope"}})
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("bad login status %d", res.StatusCode)
	}
	res.Body.Close()

	// login, follow redirect to the next page
	res, _ = cl.PostForm(srv.URL+"/admin/login", url.Values{"user": {"admin"}, "password": {testPass}, "next": {"/admin/team"}})
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.Request.URL.Path != "/admin/team" || !bytes.Contains(body, []byte("Команда")) {
		t.Fatalf("after login: %s", res.Request.URL.Path)
	}
	csrf := regexp.MustCompile(`name="_csrf" value="([^"]+)"`).FindSubmatch(body)[1]

	// POST without CSRF is refused
	res, _ = cl.PostForm(srv.URL+"/admin/contacts", url.Values{"email": {"x@y.kz"}})
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("no-csrf POST status %d", res.StatusCode)
	}
	res.Body.Close()

	// create a person with a photo
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for k, v := range map[string]string{"_csrf": string(csrf), "first_ru": "Тест", "last_ru": "Тестов", "role_perf": "1"} {
		mw.WriteField(k, v)
	}
	fw, _ := mw.CreateFormFile("photo", "p.png")
	img := image.NewRGBA(image.Rect(0, 0, 40, 50))
	img.Set(1, 1, color.White)
	png.Encode(fw, img)
	mw.Close()
	res, err := cl.Post(srv.URL+"/admin/team/new", mw.FormDataContentType(), &buf)
	if err != nil || res.Request.URL.Path != "/admin/team" {
		t.Fatalf("create person: %v %v", err, res.Status)
	}
	res.Body.Close()
	team := app.store.Get().Team
	p := team[len(team)-1]
	if p.Name.Ru != [2]string{"Тест", "Тестов"} || p.Name.En != p.Name.Ru || !strings.HasPrefix(p.Img, "/media/") {
		t.Fatalf("created person = %+v", p)
	}
	res, _ = cl.Get(srv.URL + p.Img)
	if res.StatusCode != 200 || res.Header.Get("Content-Type") != "image/jpeg" {
		t.Fatalf("uploaded photo: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	res.Body.Close()
	res, _ = cl.Get(srv.URL + "/assets/data.js")
	js, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if !bytes.Contains(js, []byte("Тестов")) {
		t.Fatal("new person not in data.js")
	}

	// delete it
	res, _ = cl.PostForm(srv.URL+"/admin/team/"+p.ID+"/delete", url.Values{"_csrf": {string(csrf)}})
	res.Body.Close()
	if len(app.store.Get().Team) != len(team)-1 {
		t.Fatal("person not deleted")
	}

	// invalid contacts are rejected and nothing changes
	old := app.store.Get().Contacts
	res, _ = cl.PostForm(srv.URL+"/admin/contacts", url.Values{"_csrf": {string(csrf)}, "phones": {"call me"}, "email": {"a@b.kz"}})
	res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity || !reflect.DeepEqual(old, app.store.Get().Contacts) {
		t.Fatalf("bad phone accepted: %d", res.StatusCode)
	}

	// a new film gets a readable id and keeps year order
	res, _ = cl.PostForm(srv.URL+"/admin/films/napoleon", url.Values{"_csrf": {string(csrf)}, "t_ru": {"Наполеон"}, "t_en": {"Napoleon"}, "year": {"2023"}, "featured": {"1"}})
	res.Body.Close()
	if res.Request.URL.Path != "/admin/films" {
		t.Fatalf("edit film landed on %s", res.Request.URL.Path)
	}
}

func TestLoginRateLimit(t *testing.T) {
	_, srv := newTestApp(t)
	var last int
	for i := 0; i < loginMaxFails+1; i++ {
		res, _ := http.PostForm(srv.URL+"/admin/login", url.Values{"user": {"admin"}, "password": {"wrong"}})
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		last = res.StatusCode
		if i == loginMaxFails && !bytes.Contains(body, []byte("Слишком много попыток")) {
			t.Fatal("no lockout after too many failures")
		}
	}
	_ = last
}

func TestOrientRotatesPhonePhotos(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 3, 2)) // 3 wide, 2 tall; mark top-left red
	src.Set(0, 0, color.RGBA{255, 0, 0, 255})
	out := orient(src, 6).(*image.RGBA) // stored sideways: display = rotate 90° clockwise
	if b := out.Bounds(); b.Dx() != 2 || b.Dy() != 3 {
		t.Fatalf("size %v", b)
	}
	if out.RGBAAt(1, 0).R != 255 { // top-left goes to top-right
		t.Fatal("orientation 6 mapped wrong")
	}
	if exifOrientation(fakeExifJPEG(6)) != 6 || exifOrientation(fakeExifJPEG(3)) != 3 {
		t.Fatal("EXIF orientation not read")
	}
}

// fakeExifJPEG builds SOI + APP1(Exif, big-endian TIFF, one orientation entry) + EOI.
func fakeExifJPEG(o uint16) []byte {
	tiff := []byte{'M', 'M', 0, 42, 0, 0, 0, 8, 0, 1, 0x01, 0x12, 0, 3, 0, 0, 0, 1, byte(o >> 8), byte(o), 0, 0, 0, 0, 0, 0, 0, 0}
	seg := append([]byte("Exif\x00\x00"), tiff...)
	n := len(seg) + 2
	b := []byte{0xFF, 0xD8, 0xFF, 0xE1, byte(n >> 8), byte(n)}
	return append(append(b, seg...), 0xFF, 0xD9)
}

func TestHelpers(t *testing.T) {
	for in, want := range map[string]string{
		"p89l3BQwIl8": "p89l3BQwIl8", "https://youtu.be/p89l3BQwIl8": "p89l3BQwIl8",
		"https://www.youtube.com/watch?v=p89l3BQwIl8&t=3": "p89l3BQwIl8", "https://example.com": "", "": "",
	} {
		if got := youtubeID(in); got != want {
			t.Errorf("youtubeID(%q) = %q", in, got)
		}
	}
	for in, want := range map[string]string{"The Way Back": "the-way-back", "Жаужүрек мың бала": "zhauzhurek-myn-bala", "  !! ": ""} {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q", in, got)
		}
	}
	if id := uniqueFilmID([]Film{{ID: "napoleon"}}, Film{Title: L10n{En: "Napoleon"}}); id != "napoleon-2" {
		t.Errorf("unique id = %s", id)
	}
}

func TestMain(m *testing.M) {
	if _, err := os.Stat(filepath.Join("..", "assets", "data.js")); err != nil {
		panic("run tests from server/: " + err.Error())
	}
	os.Exit(m.Run())
}

func get(t *testing.T, u string) (int, string) {
	t.Helper()
	res, err := http.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func TestSEOPages(t *testing.T) {
	app, srv := newTestApp(t)
	app.cfg.SiteURL = "https://nomadstunts.com"
	reLD := regexp.MustCompile(`(?s)<script type="application/ld\+json">(.*?)</script>`)
	for _, page := range []string{"/", "/about", "/team", "/projects", "/gallery", "/contacts"} {
		for _, lang := range []string{"en", "ru", "kk"} {
			u := page
			if lang != "en" {
				u += "?lang=" + lang
			}
			code, body := get(t, srv.URL+u)
			if code != 200 {
				t.Fatalf("%s: %d", u, code)
			}
			for _, want := range []string{`<html lang="` + lang + `">`, `<link rel="canonical" href="https://nomadstunts.com` + u + `">`,
				`hreflang="x-default"`, `og:image`, `<div class="ssr">`, "<h1>"} {
				if !strings.Contains(body, want) {
					t.Errorf("%s: missing %s", u, want)
				}
			}
			if strings.Count(body, `name="description"`) != 1 || strings.Count(body, "<title>") != 1 {
				t.Errorf("%s: duplicate title/description", u)
			}
			m := reLD.FindStringSubmatch(body)
			var ld map[string]any
			if m == nil || json.Unmarshal([]byte(m[1]), &ld) != nil {
				t.Fatalf("%s: JSON-LD missing or invalid", u)
			}
		}
	}
	_, about := get(t, srv.URL+"/about?lang=ru")
	for _, want := range []string{"FAQPage", "Кто основал Nomad Stunts?", "Всего проектов — 81", "Алматинская область", "+7 701 715 26 44", "<title>О Nomad Stunts"} {
		if !strings.Contains(about, want) {
			t.Errorf("about ru: missing %q", want)
		}
	}
	if strings.Contains(about, "{phones}") || strings.Contains(about, "#ZgotmplZ") {
		t.Error("about: unfilled placeholder or unsafe URL in output")
	}
	_, team := get(t, srv.URL+"/team")
	if !strings.Contains(team, "Zhandos Aibassov") || !strings.Contains(team, `"@type":"ItemList"`) {
		t.Error("team: members not server-rendered")
	}
	_, contacts := get(t, srv.URL+"/contacts")
	if !strings.Contains(contacts, `href="tel:&#43;77017152644"`) { // html/template encodes + as &#43; (valid HTML)
		t.Error("contacts: tel link not rendered")
	}
	// one URL per page
	cl := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	for from, to := range map[string]string{"/about.html?lang=ru": "/about?lang=ru", "/index.html": "/", "/index": "/"} {
		res, _ := cl.Get(srv.URL + from)
		res.Body.Close()
		if res.StatusCode != http.StatusMovedPermanently || res.Header.Get("Location") != to {
			t.Errorf("%s -> %d %s", from, res.StatusCode, res.Header.Get("Location"))
		}
	}
}

func TestCrawlerFiles(t *testing.T) {
	app, srv := newTestApp(t)
	app.cfg.SiteURL = "https://nomadstunts.com"
	_, robots := get(t, srv.URL+"/robots.txt")
	for _, want := range []string{"Disallow: /admin/", "User-agent: OAI-SearchBot", "User-agent: ClaudeBot", "Sitemap: https://nomadstunts.com/sitemap.xml"} {
		if !strings.Contains(robots, want) {
			t.Errorf("robots.txt: missing %q", want)
		}
	}
	_, sm := get(t, srv.URL+"/sitemap.xml")
	if n := strings.Count(sm, "<url>"); n != 18 {
		t.Errorf("sitemap urls = %d, want 18 (6 pages x 3 languages)", n)
	}
	if !strings.Contains(sm, "https://nomadstunts.com/projects?lang=kk") {
		t.Error("sitemap: kk projects url missing")
	}
	_, llms := get(t, srv.URL+"/llms.txt")
	for _, want := range []string{"# Nomad Stunts", "Zhaidarbek Kunguzhinov", "## FAQ", "https://nomadstunts.com/llms-full.txt"} {
		if !strings.Contains(llms, want) {
			t.Errorf("llms.txt: missing %q", want)
		}
	}
	_, full := get(t, srv.URL+"/llms-full.txt")
	if !strings.Contains(full, "## Projects") || !strings.Contains(full, "Napoleon (2023)") {
		t.Error("llms-full.txt: projects missing")
	}
}
