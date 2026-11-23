package main

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var errNotFound = errors.New("запись не найдена — возможно, её уже удалили")

// ---------- login ----------

func (a *App) loginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.auth.session(r); ok {
		http.Redirect(w, r, "/admin/", http.StatusFound)
		return
	}
	a.render(w, r, "login", Page{Title: "Вход", Data: r.URL.Query().Get("next")})
}

func (a *App) loginSubmit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
	next := r.FormValue("next")
	if !a.auth.allowLogin(r) {
		a.render(w, r, "login", Page{Title: "Вход", Error: "Слишком много попыток. Подождите 15 минут.", Data: next})
		return
	}
	if !a.auth.check(r.FormValue("user"), r.FormValue("password")) {
		a.auth.recordFail(r)
		a.render(w, r, "login", Page{Title: "Вход", Error: "Неверный логин или пароль.", Data: next})
		return
	}
	a.auth.login(w, r)
	if !strings.HasPrefix(next, "/admin/") || strings.HasPrefix(next, "//") {
		next = "/admin/"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (a *App) logoutSubmit(w http.ResponseWriter, r *http.Request) {
	a.auth.logout(w, r)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// ---------- dashboard ----------

func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	c := a.store.Get()
	a.render(w, r, "dashboard", Page{Title: "Админка", Active: "home", Data: c})
}

// ---------- showreel ----------

func (a *App) showreelForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "showreel", Page{Title: "Шоурил", Active: "showreel", Data: a.store.Get().Showreel})
}

func (a *App) showreelSave(w http.ResponseWriter, r *http.Request) {
	fail := func(msg string) {
		a.render(w, r, "showreel", Page{Title: "Шоурил", Active: "showreel", Error: msg, Data: a.store.Get().Showreel})
	}
	yt := youtubeID(r.FormValue("youtube"))
	if strings.TrimSpace(r.FormValue("youtube")) != "" && yt == "" {
		fail("Не похоже на ссылку YouTube. Вставьте ссылку вида https://youtu.be/… или оставьте поле пустым.")
		return
	}
	video, err := upload(r, "video", a.media.SaveVideo)
	if err != nil {
		fail(err.Error())
		return
	}
	poster, err := upload(r, "poster", func(f io.Reader) (string, error) { return a.media.SaveImage(f, sizeWide) })
	if err != nil {
		fail(err.Error())
		return
	}
	err = a.store.Update(func(c *Content) error {
		s := &c.Showreel
		if r.FormValue("remove_video") == "1" {
			s.MP4 = ""
		}
		if video != "" {
			s.MP4 = video
		}
		if poster != "" {
			s.Poster = poster
		}
		s.YouTube = yt
		return nil
	})
	a.done(w, r, err, "/admin/showreel", fail)
}

var ytPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^([A-Za-z0-9_-]{11})$`),
	regexp.MustCompile(`(?:youtu\.be/|/embed/|/shorts/|/live/|[?&]v=)([A-Za-z0-9_-]{11})`),
}

// youtubeID accepts a bare video id or any common YouTube link.
func youtubeID(s string) string {
	s = strings.TrimSpace(s)
	for _, p := range ytPatterns {
		if m := p.FindStringSubmatch(s); m != nil {
			return m[1]
		}
	}
	return ""
}

// ---------- contacts ----------

var phoneRe = regexp.MustCompile(`^\+?[0-9 ()\-]{5,24}$`)

func (a *App) contactsForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "contacts", Page{Title: "Контакты", Active: "contacts", Data: a.store.Get().Contacts})
}

func (a *App) contactsSave(w http.ResponseWriter, r *http.Request) {
	in := Contacts{
		Email: strings.TrimSpace(r.FormValue("email")),
		Base:  formL10n(r, "base"),
		Map:   strings.TrimSpace(r.FormValue("map")),
	}
	for _, p := range strings.Split(r.FormValue("phones"), "\n") {
		if p = strings.TrimSpace(p); p != "" {
			in.Phones = append(in.Phones, p)
		}
	}
	fail := func(msg string) {
		a.render(w, r, "contacts", Page{Title: "Контакты", Active: "contacts", Error: msg, Data: in})
	}
	for _, p := range in.Phones {
		if !phoneRe.MatchString(p) {
			fail(fmt.Sprintf("Номер «%s» выглядит неправильно. Пример: +7 701 715 26 44", p))
			return
		}
	}
	if in.Email != "" {
		if addr, err := mail.ParseAddress(in.Email); err != nil || addr.Address != in.Email {
			fail("Email выглядит неправильно.")
			return
		}
	}
	if in.Map != "" && !strings.HasPrefix(in.Map, "https://") {
		fail("Ссылка на карту должна начинаться с https://")
		return
	}
	err := a.store.Update(func(c *Content) error { c.Contacts = in; return nil })
	a.done(w, r, err, "/admin/contacts", fail)
}

// ---------- founder ----------

func (a *App) founderForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "founder", Page{Title: "Основатель", Active: "founder", Data: a.store.Get().Founder})
}

func (a *App) founderSave(w http.ResponseWriter, r *http.Request) {
	in := a.store.Get().Founder
	in.Name, in.Role, in.Hon = formName(r), formL10n(r, "role"), formL10n(r, "h")
	fail := func(msg string) {
		a.render(w, r, "founder", Page{Title: "Основатель", Active: "founder", Error: msg, Data: in})
	}
	if err := fillName(&in.Name); err != nil {
		fail(err.Error())
		return
	}
	photo, err := upload(r, "photo", func(f io.Reader) (string, error) { return a.media.SaveImage(f, sizePerson) })
	if err != nil {
		fail(err.Error())
		return
	}
	err = a.store.Update(func(c *Content) error {
		if photo != "" {
			in.Photo = photo
		}
		c.Founder = in
		return nil
	})
	a.done(w, r, err, "/admin/founder", fail)
}

// ---------- team ----------

func (a *App) teamList(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "team", Page{Title: "Команда", Active: "team", Data: a.store.Get().Team})
}

func (a *App) personForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p := Person{ID: "new", Roles: []string{"perf"}}
	if id != "new" {
		i := slices.IndexFunc(a.store.Get().Team, func(x Person) bool { return x.ID == id })
		if i < 0 {
			http.NotFound(w, r)
			return
		}
		p = a.store.Get().Team[i]
	}
	a.render(w, r, "person", Page{Title: "Участник команды", Active: "team", Data: p})
}

func (a *App) personSave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	in := Person{ID: id, Name: formName(r), Hon: formL10n(r, "h")}
	for _, role := range roleNames {
		if r.FormValue("role_"+role.Key) == "1" {
			in.Roles = append(in.Roles, role.Key)
		}
	}
	if id != "new" {
		if i := slices.IndexFunc(a.store.Get().Team, func(x Person) bool { return x.ID == id }); i >= 0 {
			in.Img, in.Photo = a.store.Get().Team[i].Img, a.store.Get().Team[i].Photo
		}
	}
	fail := func(msg string) {
		a.render(w, r, "person", Page{Title: "Участник команды", Active: "team", Error: msg, Data: in})
	}
	if err := fillName(&in.Name); err != nil {
		fail(err.Error())
		return
	}
	if len(in.Roles) == 0 {
		fail("Отметьте хотя бы одну роль.")
		return
	}
	photo, err := upload(r, "photo", func(f io.Reader) (string, error) { return a.media.SaveImage(f, sizePerson) })
	if err != nil {
		fail(err.Error())
		return
	}
	if photo != "" {
		in.Img, in.Photo = photo, photo
	}
	if in.Img == "" {
		fail("Добавьте фото.")
		return
	}
	err = a.store.Update(func(c *Content) error {
		if id == "new" {
			in.ID = randHex(5)
			c.Team = append(c.Team, in)
			return nil
		}
		i := slices.IndexFunc(c.Team, func(x Person) bool { return x.ID == id })
		if i < 0 {
			return errNotFound
		}
		c.Team[i] = in
		return nil
	})
	a.done(w, r, err, "/admin/team", fail)
}

func (a *App) personDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := a.store.Update(func(c *Content) error {
		n := len(c.Team)
		c.Team = slices.DeleteFunc(c.Team, func(x Person) bool { return x.ID == id })
		if len(c.Team) == n {
			return errNotFound
		}
		return nil
	})
	a.done(w, r, err, "/admin/team", a.failText(w))
}

func (a *App) personMove(w http.ResponseWriter, r *http.Request) {
	id, up := r.PathValue("id"), r.FormValue("dir") == "up"
	err := a.store.Update(func(c *Content) error {
		return move(c.Team, func(x Person) bool { return x.ID == id }, up)
	})
	a.done(w, r, err, "/admin/team#p-"+id, a.failText(w))
}

// ---------- films ----------

var filmIDRe = regexp.MustCompile(`^[a-z0-9_-]+$`)

func (a *App) filmList(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "films", Page{Title: "Проекты", Active: "films", Data: a.store.Get().Films})
}

func (a *App) filmForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f := Film{ID: "new"}
	if id != "new" {
		films := a.store.Get().Films
		i := slices.IndexFunc(films, func(x Film) bool { return x.ID == id })
		if i < 0 {
			http.NotFound(w, r)
			return
		}
		f = films[i]
	}
	a.render(w, r, "film", Page{Title: "Проект", Active: "films", Data: f})
}

func (a *App) filmSave(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id != "new" && !filmIDRe.MatchString(id) {
		http.NotFound(w, r)
		return
	}
	in := Film{
		ID:    id,
		Title: formL10n(r, "t"),
		En:    formInfo(r, "en"),
		Ru:    formInfo(r, "ru"),
		Kk:    formInfo(r, "kk"),
	}
	in.Year, _ = strconv.Atoi(strings.TrimSpace(r.FormValue("year")))
	in.Featured, _ = strconv.Atoi(strings.TrimSpace(r.FormValue("featured")))
	if id != "new" {
		films := a.store.Get().Films
		if i := slices.IndexFunc(films, func(x Film) bool { return x.ID == id }); i >= 0 {
			in.Poster, in.Img = films[i].Poster, films[i].Img
		}
	}
	fail := func(msg string) {
		a.render(w, r, "film", Page{Title: "Проект", Active: "films", Error: msg, Data: in})
	}
	switch {
	case in.Title.Ru == "" && in.Title.En == "":
		fail("Укажите название хотя бы на русском или английском.")
		return
	case in.Year < 1950 || in.Year > 2100:
		fail("Укажите год выхода, например 2024.")
		return
	case in.Featured < 0 || in.Featured > 999:
		fail("Позиция на главной — число от 1 до 999, или оставьте пустым.")
		return
	}
	poster, err := upload(r, "poster", func(f io.Reader) (string, error) { return a.media.SaveImage(f, sizePoster) })
	if err != nil {
		fail(err.Error())
		return
	}
	if poster != "" {
		in.Poster, in.Img = poster, poster
	}
	if in.Poster == "" {
		fail("Добавьте постер.")
		return
	}
	err = a.store.Update(func(c *Content) error {
		if id == "new" {
			in.ID = uniqueFilmID(c.Films, in)
			c.Films = append(c.Films, in)
		} else {
			i := slices.IndexFunc(c.Films, func(x Film) bool { return x.ID == id })
			if i < 0 {
				return errNotFound
			}
			c.Films[i] = in
		}
		sort.SliceStable(c.Films, func(i, j int) bool { return c.Films[i].Year > c.Films[j].Year })
		return nil
	})
	a.done(w, r, err, "/admin/films", fail)
}

func (a *App) filmDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := a.store.Update(func(c *Content) error {
		n := len(c.Films)
		c.Films = slices.DeleteFunc(c.Films, func(x Film) bool { return x.ID == id })
		if len(c.Films) == n {
			return errNotFound
		}
		return nil
	})
	a.done(w, r, err, "/admin/films", a.failText(w))
}

// uniqueFilmID makes a readable URL id (#napoleon) from the title.
func uniqueFilmID(films []Film, f Film) string {
	base := slugify(f.Title.En)
	if base == "" {
		base = slugify(f.Title.Ru)
	}
	if base == "" {
		base = "film-" + strconv.Itoa(f.Year)
	}
	id := base
	for n := 2; id == "new" || slices.ContainsFunc(films, func(x Film) bool { return x.ID == id }); n++ { // "new" is the create route
		id = base + "-" + strconv.Itoa(n)
	}
	return id
}

var translit = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i", 'й': "y",
	'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f",
	'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	'ә': "a", 'ғ': "g", 'қ': "k", 'ң': "n", 'ө': "o", 'ұ': "u", 'ү': "u", 'һ': "h", 'і': "i",
}

func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r < unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)):
			b.WriteRune(r)
			dash = false
		case translit[r] != "":
			b.WriteString(translit[r])
			dash = false
		case b.Len() > 0 && !dash:
			b.WriteByte('-')
			dash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 60 {
		out = strings.Trim(out[:60], "-")
	}
	return out
}

// ---------- gallery ----------

const maxGalleryFiles = 40

func (a *App) galleryPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "gallery", Page{Title: "Фотогалерея", Active: "gallery", Data: a.store.Get().Gallery})
}

func (a *App) galleryUpload(w http.ResponseWriter, r *http.Request) {
	fail := func(msg string) {
		a.render(w, r, "gallery", Page{Title: "Фотогалерея", Active: "gallery", Error: msg, Data: a.store.Get().Gallery})
	}
	var files []*multipart.FileHeader
	if r.MultipartForm != nil {
		files = r.MultipartForm.File["photos"]
	}
	if len(files) == 0 {
		fail("Выберите одно или несколько фото.")
		return
	}
	if len(files) > maxGalleryFiles {
		fail(fmt.Sprintf("За один раз можно загрузить до %d фото.", maxGalleryFiles))
		return
	}
	var added []Photo
	var bad []string
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			bad = append(bad, fh.Filename)
			continue
		}
		full, thumb, err := a.media.SaveImageWithThumb(f, sizeGallery, sizeThumb)
		f.Close()
		if err != nil {
			bad = append(bad, fh.Filename)
			continue
		}
		added = append(added, Photo{ID: randHex(5), Img: full, Thumb: thumb})
	}
	if len(added) > 0 {
		err := a.store.Update(func(c *Content) error {
			c.Gallery = append(added, c.Gallery...) // newest first
			return nil
		})
		if err != nil {
			fail(err.Error())
			return
		}
	}
	if len(bad) > 0 {
		a.render(w, r, "gallery", Page{Title: "Фотогалерея", Active: "gallery", Data: a.store.Get().Gallery,
			Flash: fmt.Sprintf("Добавлено фото: %d.", len(added)),
			Error: "Не удалось загрузить (это не JPEG/PNG/WebP или файл повреждён): " + strings.Join(bad, ", ")})
		return
	}
	http.Redirect(w, r, "/admin/gallery?saved=1", http.StatusSeeOther)
}

func (a *App) photoDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := a.store.Update(func(c *Content) error {
		n := len(c.Gallery)
		c.Gallery = slices.DeleteFunc(c.Gallery, func(x Photo) bool { return x.ID == id })
		if len(c.Gallery) == n {
			return errNotFound
		}
		return nil
	})
	a.done(w, r, err, "/admin/gallery", a.failText(w))
}

func (a *App) photoMove(w http.ResponseWriter, r *http.Request) {
	id, up := r.PathValue("id"), r.FormValue("dir") == "up"
	err := a.store.Update(func(c *Content) error {
		return move(c.Gallery, func(x Photo) bool { return x.ID == id }, up)
	})
	a.done(w, r, err, "/admin/gallery#g-"+id, a.failText(w))
}

// ---------- history ----------

func (a *App) historyPage(w http.ResponseWriter, r *http.Request) {
	snaps, err := a.store.History()
	p := Page{Title: "История изменений", Active: "history", Data: snaps}
	if err != nil {
		p.Error = err.Error()
	}
	a.render(w, r, "history", p)
}

func (a *App) historyRestore(w http.ResponseWriter, r *http.Request) {
	err := a.store.Restore(r.FormValue("name"))
	a.done(w, r, err, "/admin/history", a.failText(w))
}

// ---------- helpers ----------

// done redirects after a successful save (PRG), or shows the error on the form.
func (a *App) done(w http.ResponseWriter, r *http.Request, err error, to string, fail func(string)) {
	if err != nil {
		fail(err.Error())
		return
	}
	target, frag, _ := strings.Cut(to, "#")
	if frag != "" {
		frag = "#" + frag
	}
	http.Redirect(w, r, target+"?saved=1"+frag, http.StatusSeeOther)
}

func (a *App) failText(w http.ResponseWriter) func(string) {
	return func(msg string) { http.Error(w, msg, http.StatusUnprocessableEntity) }
}

// upload saves an optional file field; ("", nil) when no file was chosen.
func upload(r *http.Request, field string, save func(io.Reader) (string, error)) (string, error) {
	f, h, err := r.FormFile(field)
	if errors.Is(err, http.ErrMissingFile) || errors.Is(err, http.ErrNotMultipart) {
		return "", nil
	}
	if err != nil {
		return "", errors.New("не удалось прочитать файл")
	}
	defer f.Close()
	if h.Size == 0 {
		return "", nil
	}
	return save(f)
}

func move[T any](list []T, match func(T) bool, up bool) error {
	i := slices.IndexFunc(list, match)
	if i < 0 {
		return errNotFound
	}
	j := i + 1
	if up {
		j = i - 1
	}
	if j >= 0 && j < len(list) {
		list[i], list[j] = list[j], list[i]
	}
	return nil
}

func formL10n(r *http.Request, prefix string) L10n {
	v := func(l string) string { return strings.TrimSpace(r.FormValue(prefix + "_" + l)) }
	return L10n{En: v("en"), Ru: v("ru"), Kk: v("kk")}
}

func formInfo(r *http.Request, lang string) FilmInfo {
	v := func(k string) string { return strings.TrimSpace(r.FormValue(k + "_" + lang)) }
	return FilmInfo{Type: v("type"), Country: v("country"), Director: v("director"), Stars: v("stars")}
}

func formName(r *http.Request) Name {
	v := func(k string) string { return strings.TrimSpace(r.FormValue(k)) }
	return Name{
		En: [2]string{v("first_en"), v("last_en")},
		Ru: [2]string{v("first_ru"), v("last_ru")},
		Kk: [2]string{v("first_kk"), v("last_kk")},
	}
}

// fillName requires a name and copies it into languages left empty,
// because the site would otherwise show a blank name in that language.
func fillName(n *Name) error {
	empty := func(x [2]string) bool { return x[0] == "" && x[1] == "" }
	src := n.Ru
	if empty(src) {
		src = n.En
	}
	if empty(src) {
		return errors.New("Укажите имя и фамилию хотя бы на русском.")
	}
	for _, x := range []*[2]string{&n.Ru, &n.Kk, &n.En} {
		if empty(*x) {
			*x = src
		}
	}
	return nil
}
