package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// SEO: every page is served with per-language <title>, description, canonical, hreflang,
// Open Graph and JSON-LD, plus a plain-HTML copy of its content (.ssr) for crawlers and
// AI bots that don't run JavaScript. The JS page hides that copy and renders the real UI.

//go:embed seo.html
var seoFS embed.FS

var langs = []string{"en", "ru", "kk"} // en is the default (no ?lang)

var pageKeys = map[string]string{"index": "home", "about": "about", "team": "team", "projects": "projects", "gallery": "gallery", "contacts": "contacts"}
var pageOrder = []string{"home", "about", "team", "projects", "gallery", "contacts"}

var social = []struct{ Name, URL string }{
	{"Instagram", "https://www.instagram.com/nomadstunts/"},
	{"Facebook", "https://www.facebook.com/nomadstunts"},
	{"YouTube", "https://www.youtube.com/channel/UCx2yGDE7QSX8n6cWHRJ2i7g"},
}

// I18N mirrors assets/i18n.js (window.NS_I18N), the texts shared with the site JS.
type I18N struct {
	T          map[string]L10n                 `json:"t"`
	Services   []Service                       `json:"services"`
	Directions [][]string                      `json:"directions"` // [icon, en, ru, kk]
	History    map[string][][2]string          `json:"history"`
	Stars      map[string][]string             `json:"stars"`
	FAQ        map[string][][2]string          `json:"faq"`
	SEO        map[string]map[string][2]string `json:"seo"` // page -> lang -> [title, description]
}

type Service struct {
	Img string    `json:"img"`
	En  [2]string `json:"en"`
	Ru  [2]string `json:"ru"`
	Kk  [2]string `json:"kk"`
}

func loadI18N(siteDir string) (*I18N, error) {
	raw, err := readDataJS(filepath.Join(siteDir, "assets", "i18n.js"))
	if err != nil {
		return nil, fmt.Errorf("i18n.js: %w", err)
	}
	var in I18N
	if err := decodeStrict(raw, &in); err != nil {
		return nil, fmt.Errorf("i18n.js: %w", err)
	}
	return &in, nil
}

func (l L10n) In(lang string) string {
	v := map[string]string{"en": l.En, "ru": l.Ru, "kk": l.Kk}[lang]
	if v == "" {
		v = l.Ru
	}
	if v == "" {
		v = l.En
	}
	return v
}

func (n Name) In(lang string) string {
	x := map[string][2]string{"en": n.En, "ru": n.Ru, "kk": n.Kk}[lang]
	return strings.TrimSpace(x[0] + " " + x[1])
}

func (s Service) In(lang string) [2]string {
	return map[string][2]string{"en": s.En, "ru": s.Ru, "kk": s.Kk}[lang]
}

func (f Film) InfoIn(lang string) FilmInfo {
	for _, l := range []string{lang, "en", "ru"} {
		if i := map[string]FilmInfo{"en": f.En, "ru": f.Ru, "kk": f.Kk}[l]; i != (FilmInfo{}) {
			return i
		}
	}
	return FilmInfo{}
}

func (f Film) TitleIn(lang string) string {
	if t := map[string]string{"en": f.Title.En, "ru": f.Title.Ru, "kk": f.Title.Kk}[lang]; t != "" {
		return t
	}
	if f.Title.En != "" {
		return f.Title.En
	}
	return f.Title.Ru
}

// ---------- page view ----------

type link struct {
	Href, Label, Lang string
	On                bool
}

type seoPage struct {
	Lang, Page, Base, Canonical, Title, Desc, OGLocale string
	Alternates                                         []link
	OGAlt                                              []string
	JSONLD                                             template.JS
	Nav, Langs                                         []link
	C                                                  Content
	i                                                  *I18N
}

// T returns a developer-authored interface text (may contain markup).
func (p *seoPage) T(key string) template.HTML { return template.HTML(p.i.T[key].In(p.Lang)) }
func (p *seoPage) Plain(key string) string    { return stripTags(p.i.T[key].In(p.Lang)) }
func (p *seoPage) Abs(path string) string     { return absURL(p.Base, path) }
func (p *seoPage) Founder() string            { return p.C.Founder.Name.In(p.Lang) }
func (p *seoPage) Roles(roles []string) string {
	var out []string
	for _, r := range roles {
		out = append(out, p.i.T["r_"+r].In(p.Lang))
	}
	return strings.Join(out, ", ")
}
func (p *seoPage) Services() [][2]string {
	var out [][2]string
	for _, s := range p.i.Services {
		out = append(out, s.In(p.Lang))
	}
	return out
}
func (p *seoPage) Directions() []string {
	idx := map[string]int{"en": 1, "ru": 2, "kk": 3}[p.Lang]
	var out []string
	for _, d := range p.i.Directions {
		if idx < len(d) {
			out = append(out, d[idx])
		}
	}
	return out
}
func (p *seoPage) History() [][2]string { return p.i.History[p.Lang] }
func (p *seoPage) FAQ() [][2]string     { return faqIn(p.i, p.C, p.Lang) }
func (p *seoPage) Featured() []Film     { return featured(p.C.Films) }
func (p *seoPage) TopTeam() []Person    { return p.C.Team[:min(12, len(p.C.Team))] }
func (p *seoPage) Stars() string        { return strings.Join(p.i.Stars[p.Lang], ", ") }
func (p *seoPage) Social() any          { return social }
func (p *seoPage) Tel(phone string) template.URL {
	return template.URL("tel:" + strings.Map(func(r rune) rune {
		if r == '+' || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, phone))
}

var seoTmpl = template.Must(template.New("seo.html").Funcs(template.FuncMap{
	"html2": func(s string) template.HTML { return template.HTML(s) }, // developer-authored texts from i18n.js
	"inc":   func(i int) int { return i + 1 },
}).ParseFS(seoFS, "seo.html"))

var (
	reHTMLLang = regexp.MustCompile(`<html lang="[^"]*">`)
	reTitle    = regexp.MustCompile(`<title>[^<]*</title>`)
	reDesc     = regexp.MustCompile(`<meta name="description" content="[^"]*">\n?`)
	reBody     = regexp.MustCompile(`<body[^>]*>`)
)

// servePage renders SITE_DIR/<name>.html for the requested language.
func (a *App) servePage(w http.ResponseWriter, r *http.Request, name string) {
	raw, err := os.ReadFile(filepath.Join(a.cfg.SiteDir, name+".html"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	p := a.seoPage(r, pageKeys[name], pickLang(r))
	var head, ssr bytes.Buffer
	if err := seoTmpl.ExecuteTemplate(&head, "head", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := seoTmpl.ExecuteTemplate(&ssr, "ssr", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := reHTMLLang.ReplaceAll(raw, []byte(`<html lang="`+p.Lang+`">`))
	out = reTitle.ReplaceAllLiteral(out, []byte("<title>"+template.HTMLEscapeString(p.Title)+"</title>"))
	out = reDesc.ReplaceAllLiteral(out, nil)
	out = bytes.Replace(out, []byte("</head>"), append(head.Bytes(), "</head>"...), 1)
	if loc := reBody.FindIndex(out); loc != nil {
		out = append(out[:loc[1]:loc[1]], append(append([]byte("\n"), ssr.Bytes()...), out[loc[1]:]...)...)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Language", p.Lang)
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(out)
}

func pickLang(r *http.Request) string {
	l := r.URL.Query().Get("lang")
	for _, x := range langs {
		if l == x {
			return l
		}
	}
	return "en"
}

func (a *App) baseURL(r *http.Request) string {
	if a.cfg.SiteURL != "" {
		return a.cfg.SiteURL
	}
	scheme := "http"
	if isHTTPS(r) {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func pageURL(base, page, lang string) string {
	p := "/"
	if page != "home" {
		p += page
	}
	if lang != "en" {
		p += "?lang=" + lang
	}
	return base + p
}

func (a *App) seoPage(r *http.Request, page, lang string) *seoPage {
	base := a.baseURL(r)
	c := a.store.Get()
	meta := a.i18n.SEO[page][lang]
	if meta[0] == "" {
		meta = a.i18n.SEO[page]["en"]
	}
	locales := map[string]string{"en": "en_US", "ru": "ru_RU", "kk": "kk_KZ"}
	p := &seoPage{
		Lang: lang, Page: page, Base: base, Canonical: pageURL(base, page, lang),
		Title: meta[0], Desc: meta[1], OGLocale: locales[lang], C: c, i: a.i18n,
	}
	for _, l := range langs {
		p.Alternates = append(p.Alternates, link{Lang: l, Href: pageURL(base, page, l)})
		if l != lang {
			p.OGAlt = append(p.OGAlt, locales[l])
			p.Langs = append(p.Langs, link{Lang: l, Href: pageURL(base, page, l), Label: map[string]string{"en": "English", "ru": "Русский", "kk": "Қазақша"}[l]})
		}
	}
	for _, pg := range pageOrder {
		key := map[string]string{"home": "home", "about": "nav_about", "team": "nav_team", "projects": "nav_projects", "gallery": "nav_gallery", "contacts": "nav_contacts"}[pg]
		p.Nav = append(p.Nav, link{Href: pageURL(base, pg, lang), Label: a.i18n.T[key].In(lang), On: pg == page})
	}
	b, _ := json.Marshal(a.jsonLD(p))
	p.JSONLD = template.JS(b)
	return p
}

// ---------- structured data ----------

func (a *App) jsonLD(p *seoPage) map[string]any {
	c, base, lang := p.C, p.Base, p.Lang
	org := map[string]any{
		"@type": "Organization", "@id": base + "/#org", "name": "Nomad Stunts",
		"alternateName": "Nomad Stunts International Action Company",
		"url":           base + "/", "logo": base + "/assets/logo.svg", "image": base + "/assets/og-image.jpg",
		"description": stripTags(a.i18n.T["intro_p"].In(lang)), "foundingDate": "2000",
		"founder": map[string]any{"@id": base + "/#founder"}, "areaServed": "Worldwide",
		"address":    map[string]any{"@type": "PostalAddress", "addressLocality": "Issyk", "addressRegion": "Almaty Region", "addressCountry": "KZ"},
		"knowsAbout": servicesTitles(a.i18n, lang),
	}
	var sameAs []string
	for _, s := range social {
		sameAs = append(sameAs, s.URL)
	}
	org["sameAs"] = sameAs
	if c.Contacts.Email != "" {
		org["email"] = c.Contacts.Email
	}
	if len(c.Contacts.Phones) > 0 {
		org["telephone"] = c.Contacts.Phones[0]
		cp := map[string]any{"@type": "ContactPoint", "contactType": "booking", "telephone": c.Contacts.Phones, "availableLanguage": []string{"Kazakh", "Russian", "English"}}
		if c.Contacts.Email != "" {
			cp["email"] = c.Contacts.Email
		}
		org["contactPoint"] = cp
	}
	if c.Contacts.Map != "" {
		org["hasMap"] = c.Contacts.Map
	}
	founder := map[string]any{
		"@type": "Person", "@id": base + "/#founder", "name": c.Founder.Name.In(lang),
		"alternateName": distinct(c.Founder.Name.In("en"), c.Founder.Name.In("ru"), c.Founder.Name.In("kk")),
		"description":   c.Founder.Role.In(lang), "worksFor": map[string]any{"@id": base + "/#org"},
		"jobTitle": a.i18n.T["r_coord"].In(lang),
		"award":    []string{"Taurus World Stunt Awards 2020", stripTags(a.i18n.T["aw2"].In(lang)), stripTags(a.i18n.T["aw3"].In(lang))},
		"memberOf": map[string]any{"@type": "Organization", "name": "Hitz International", "url": "https://www.hitzinternational.com/"},
	}
	if c.Founder.Photo != "" {
		founder["image"] = absURL(base, c.Founder.Photo)
	}
	site := map[string]any{"@type": "WebSite", "@id": base + "/#website", "url": base + "/", "name": "Nomad Stunts", "inLanguage": langs, "publisher": map[string]any{"@id": base + "/#org"}}
	pageType := map[string]string{"home": "WebPage", "about": "AboutPage", "team": "CollectionPage", "projects": "CollectionPage", "gallery": "ImageGallery", "contacts": "ContactPage"}[p.Page]
	webpage := map[string]any{
		"@type": pageType, "@id": p.Canonical + "#webpage", "url": p.Canonical, "name": p.Title, "description": p.Desc,
		"inLanguage": lang, "isPartOf": map[string]any{"@id": base + "/#website"}, "about": map[string]any{"@id": base + "/#org"},
	}
	graph := []any{org, founder, site, webpage}
	if p.Page != "home" {
		webpage["breadcrumb"] = map[string]any{"@type": "BreadcrumbList", "itemListElement": []any{
			map[string]any{"@type": "ListItem", "position": 1, "name": a.i18n.T["home"].In(lang), "item": pageURL(base, "home", lang)},
			map[string]any{"@type": "ListItem", "position": 2, "name": p.Nav[indexOf(pageOrder, p.Page)].Label, "item": p.Canonical},
		}}
	}
	switch p.Page {
	case "about":
		var qs []any
		for _, qa := range faqIn(a.i18n, c, lang) {
			qs = append(qs, map[string]any{"@type": "Question", "name": qa[0], "acceptedAnswer": map[string]any{"@type": "Answer", "text": qa[1]}})
		}
		graph = append(graph, map[string]any{"@type": "FAQPage", "@id": p.Canonical + "#faq", "inLanguage": lang, "mainEntity": qs})
	case "team":
		var items []any
		for i, m := range c.Team {
			person := map[string]any{"@type": "Person", "name": m.Name.In(lang), "jobTitle": p.Roles(m.Roles), "memberOf": map[string]any{"@id": base + "/#org"}}
			if m.Img != "" {
				person["image"] = absURL(base, m.Img)
			}
			items = append(items, map[string]any{"@type": "ListItem", "position": i + 1, "item": person})
		}
		webpage["mainEntity"] = map[string]any{"@type": "ItemList", "numberOfItems": len(items), "itemListElement": items}
	case "projects":
		var items []any
		for i, f := range c.Films {
			info := f.InfoIn(lang)
			typ := "Movie"
			if strings.Contains(strings.ToLower(f.InfoIn("en").Type+f.InfoIn("ru").Type), "seri") || strings.Contains(strings.ToLower(f.InfoIn("ru").Type), "сериал") {
				typ = "TVSeries"
			}
			work := map[string]any{"@type": typ, "name": f.TitleIn(lang), "dateCreated": strconv.Itoa(f.Year), "contributor": map[string]any{"@id": base + "/#org"}}
			if f.Poster != "" {
				work["image"] = absURL(base, f.Poster)
			}
			if info.Director != "" {
				work["director"] = map[string]any{"@type": "Person", "name": info.Director}
			}
			if info.Country != "" {
				work["countryOfOrigin"] = map[string]any{"@type": "Country", "name": info.Country}
			}
			items = append(items, map[string]any{"@type": "ListItem", "position": i + 1, "item": work})
		}
		webpage["mainEntity"] = map[string]any{"@type": "ItemList", "numberOfItems": len(items), "itemListElement": items}
	case "gallery":
		var media []any
		for _, g := range c.Gallery {
			media = append(media, map[string]any{"@type": "ImageObject", "contentUrl": absURL(base, g.Img), "thumbnailUrl": absURL(base, g.Thumb)})
		}
		webpage["associatedMedia"] = media
	}
	return map[string]any{"@context": "https://schema.org", "@graph": graph}
}

// ---------- shared content helpers ----------

// faqIn fills {projects}, {base}, {phones}, {email} with live content (same rule as app.js).
func faqIn(i *I18N, c Content, lang string) [][2]string {
	r := strings.NewReplacer("{projects}", strconv.Itoa(len(c.Films)), "{base}", c.Contacts.Base.In(lang),
		"{phones}", strings.Join(c.Contacts.Phones, ", "), "{email}", c.Contacts.Email)
	var out [][2]string
	for _, qa := range i.FAQ[lang] {
		out = append(out, [2]string{qa[0], r.Replace(qa[1])})
	}
	return out
}

func featured(films []Film) []Film {
	var out []Film
	for n := 1; n <= 999 && len(out) < len(films); n++ {
		for _, f := range films {
			if f.Featured == n {
				out = append(out, f)
			}
		}
	}
	return out
}

func servicesTitles(i *I18N, lang string) []string {
	var out []string
	for _, s := range i.Services {
		out = append(out, s.In(lang)[0])
	}
	return out
}

var reTags = regexp.MustCompile(`<[^>]+>`)

func stripTags(s string) string { return strings.TrimSpace(reTags.ReplaceAllString(s, "")) }

func absURL(base, p string) string {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	u := strings.TrimPrefix(p, ".")
	if !strings.HasPrefix(u, "/") {
		u = "/upload/images/" + u
	}
	return base + (&url.URL{Path: u}).EscapedPath()
}

func distinct(xs ...string) []string {
	var out []string
	for _, x := range xs {
		if x != "" && indexOf(out, x) < 0 {
			out = append(out, x)
		}
	}
	return out
}

func indexOf(xs []string, v string) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return -1
}

// ---------- robots.txt, sitemap.xml, llms.txt ----------

func (a *App) robots(w http.ResponseWriter, r *http.Request) {
	base := a.baseURL(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, `# Nomad Stunts — search engines and AI assistants are welcome to read and cite this site.
User-agent: *
Allow: /
Disallow: /admin/
Content-Signal: search=yes, ai-input=yes, ai-train=yes

# AI search, user-requested fetches and model crawlers (listed explicitly; same rules as above)
User-agent: OAI-SearchBot
User-agent: ChatGPT-User
User-agent: GPTBot
User-agent: Claude-SearchBot
User-agent: Claude-User
User-agent: ClaudeBot
User-agent: PerplexityBot
User-agent: Perplexity-User
User-agent: Google-Extended
User-agent: Applebot-Extended
User-agent: Bingbot
Allow: /
Disallow: /admin/

Sitemap: %s/sitemap.xml
`, base)
}

func (a *App) sitemap(w http.ResponseWriter, r *http.Request) {
	base := a.baseURL(r)
	mod := a.store.Modified()
	for _, name := range []string{"index", "about", "team", "projects", "gallery", "contacts"} {
		if st, err := os.Stat(filepath.Join(a.cfg.SiteDir, name+".html")); err == nil && st.ModTime().After(mod) {
			mod = st.ModTime()
		}
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	for _, pg := range pageOrder {
		for _, l := range langs {
			fmt.Fprintf(&b, "<url><loc>%s</loc><lastmod>%s</lastmod>", xmlEsc(pageURL(base, pg, l)), mod.UTC().Format("2006-01-02"))
			for _, alt := range langs {
				fmt.Fprintf(&b, `<xhtml:link rel="alternate" hreflang="%s" href="%s"/>`, alt, xmlEsc(pageURL(base, pg, alt)))
			}
			fmt.Fprintf(&b, `<xhtml:link rel="alternate" hreflang="x-default" href="%s"/></url>`+"\n", xmlEsc(pageURL(base, pg, "en")))
		}
	}
	b.WriteString("</urlset>\n")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(b.String()))
}

func xmlEsc(s string) string { return template.HTMLEscapeString(s) }

// llms.txt (llmstxt.org): a plain-text brief for AI assistants, built from live content.
func (a *App) llms(full bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base, c, i := a.baseURL(r), a.store.Get(), a.i18n
		var b strings.Builder
		b.WriteString("# Nomad Stunts\n\n> " + i.SEO["home"]["en"][1] + "\n\n")
		b.WriteString(stripTags(i.T["intro_p"].En) + "\n\n## Key facts\n\n")
		var feat []string
		for _, f := range featured(c.Films) {
			feat = append(feat, f.TitleIn("en")+" ("+strconv.Itoa(f.Year)+")")
		}
		fmt.Fprintf(&b, "- Founded: 2000, Kazakhstan\n- Founder: %s — %s\n", c.Founder.Name.In("en"), c.Founder.Role.In("en"))
		fmt.Fprintf(&b, "- Team on the site: %d people (stunt coordinators, stunt performers, horse masters, riggers)\n", len(c.Team)+1)
		fmt.Fprintf(&b, "- Training base: %s\n- Projects: %d films and series, including %s\n", c.Contacts.Base.In("en"), len(c.Films), strings.Join(feat, ", "))
		fmt.Fprintf(&b, "- Worked with: %s\n- Services: %s\n", strings.Join(i.Stars["en"], ", "), strings.Join(servicesTitles(i, "en"), ", "))
		fmt.Fprintf(&b, "- Contacts: %s; %s\n\n## Pages\n\n", strings.Join(c.Contacts.Phones, ", "), c.Contacts.Email)
		descs := map[string]string{"home": "overview", "about": "history, services, training base, FAQ", "team": "founder and full team with roles",
			"projects": "all films and series by year", "gallery": "behind-the-scenes photos", "contacts": "phones, email, training base location"}
		for _, pg := range pageOrder {
			fmt.Fprintf(&b, "- [%s](%s): %s\n", i.SEO[pg]["en"][0], pageURL(base, pg, "en"), descs[pg])
		}
		b.WriteString("\nRussian and Kazakh versions: add `?lang=ru` or `?lang=kk` to any page URL.\n\n## FAQ\n\n")
		for _, qa := range faqIn(i, c, "en") {
			b.WriteString("### " + qa[0] + "\n\n" + qa[1] + "\n\n")
		}
		if !full {
			fmt.Fprintf(&b, "## Optional\n\n- [Full team, projects and history](%s/llms-full.txt)\n", base)
		} else {
			b.WriteString("## History\n\n")
			for _, h := range i.History["en"] {
				b.WriteString("- **" + h[0] + "** — " + stripTags(h[1]) + "\n")
			}
			b.WriteString("\n## Services\n\n")
			for _, s := range i.Services {
				b.WriteString("- **" + s.En[0] + "** — " + s.En[1] + "\n")
			}
			b.WriteString("\n## Team\n\n- " + c.Founder.Name.In("en") + " — founder. " + c.Founder.Hon.In("en") + "\n")
			for _, m := range c.Team {
				line := "- " + m.Name.In("en") + " — " + rolesEn(i, m.Roles)
				if h := m.Hon.In("en"); h != "" {
					line += ". " + h
				}
				b.WriteString(line + "\n")
			}
			b.WriteString("\n## Projects\n\n")
			for _, f := range c.Films {
				info := f.InfoIn("en")
				line := "- " + f.TitleIn("en") + " (" + strconv.Itoa(f.Year) + ")"
				for _, kv := range [][2]string{{"type", info.Type}, {"country", info.Country}, {"director", info.Director}, {"cast", info.Stars}} {
					if kv[1] != "" {
						line += "; " + kv[0] + ": " + kv[1]
					}
				}
				b.WriteString(line + "\n")
			}
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write([]byte(b.String()))
	}
}

func rolesEn(i *I18N, roles []string) string {
	var out []string
	for _, r := range roles {
		out = append(out, strings.ToLower(i.T["r_"+r].En))
	}
	return strings.Join(out, ", ")
}
