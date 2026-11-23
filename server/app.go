package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

type App struct {
	cfg   Config
	store *Store
	media *Media
	auth  *Auth
	pages map[string]*template.Template
	i18n  *I18N
}

type ctxKey struct{}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	static, _ := fs.Sub(staticFS, "static")

	mux.HandleFunc("GET /assets/data.js", a.dataJS)
	mux.HandleFunc("GET /robots.txt", a.robots)
	mux.HandleFunc("GET /sitemap.xml", a.sitemap)
	mux.HandleFunc("GET /llms.txt", a.llms(false))
	mux.HandleFunc("GET /llms-full.txt", a.llms(true))
	mux.Handle("GET /media/", a.files(filepath.Join(a.cfg.DataDir, "media"), "/media/", true))
	mux.Handle("GET /admin/static/", http.StripPrefix("/admin/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/admin/", http.StatusFound) })
	mux.HandleFunc("GET /admin/login", a.loginPage)
	mux.HandleFunc("POST /admin/login", a.loginSubmit)

	admin := func(pattern string, h http.HandlerFunc) { mux.Handle(pattern, a.protected(h)) }
	admin("POST /admin/logout", a.logoutSubmit)
	admin("GET /admin/{$}", a.dashboard)
	admin("GET /admin/showreel", a.showreelForm)
	admin("POST /admin/showreel", a.showreelSave)
	admin("GET /admin/contacts", a.contactsForm)
	admin("POST /admin/contacts", a.contactsSave)
	admin("GET /admin/founder", a.founderForm)
	admin("POST /admin/founder", a.founderSave)
	admin("GET /admin/team", a.teamList)
	admin("GET /admin/team/{id}", a.personForm)
	admin("POST /admin/team/{id}", a.personSave)
	admin("POST /admin/team/{id}/delete", a.personDelete)
	admin("POST /admin/team/{id}/move", a.personMove)
	admin("GET /admin/films", a.filmList)
	admin("GET /admin/films/{id}", a.filmForm)
	admin("POST /admin/films/{id}", a.filmSave)
	admin("POST /admin/films/{id}/delete", a.filmDelete)
	admin("GET /admin/gallery", a.galleryPage)
	admin("POST /admin/gallery", a.galleryUpload)
	admin("POST /admin/gallery/{id}/delete", a.photoDelete)
	admin("POST /admin/gallery/{id}/move", a.photoMove)
	admin("GET /admin/history", a.historyPage)
	admin("POST /admin/history/restore", a.historyRestore)

	mux.Handle("GET /", a.site())
	return headers(mux)
}

// protected requires a valid session; POSTs also get a body limit and a CSRF check.
func (a *App) protected(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := a.auth.session(r)
		if !ok {
			if r.Method == http.MethodGet {
				http.Redirect(w, r, "/admin/login?next="+r.URL.Path, http.StatusFound)
			} else {
				http.Error(w, "Сессия истекла. Войдите снова.", http.StatusForbidden)
			}
			return
		}
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, maxVideoBytes+16<<20)
			var err error
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				err = r.ParseMultipartForm(32 << 20)
			} else {
				err = r.ParseForm()
			}
			var tooBig *http.MaxBytesError
			if errors.As(err, &tooBig) {
				http.Error(w, "Файлы слишком большие для одной загрузки.", http.StatusRequestEntityTooLarge)
				return
			}
			if err != nil {
				http.Error(w, "Не удалось прочитать форму.", http.StatusBadRequest)
				return
			}
			if !a.auth.validCSRF(r, s) {
				http.Error(w, "Страница устарела. Обновите её и повторите.", http.StatusForbidden)
				return
			}
		}
		h(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, s)))
	})
}

func (a *App) dataJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(a.store.DataJS())
}

// site serves the static landing page: top-level *.html (also without ".html"), /assets/, /upload/.
// Nothing else in SITE_DIR (server code, .git, data) is reachable.
func (a *App) site() http.Handler {
	assets := a.files(a.cfg.SiteDir, "/", false)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean("/" + r.URL.Path)
		switch {
		case strings.Contains(p, "/."):
			http.NotFound(w, r)
		case p == "/favicon.ico" || strings.HasPrefix(p, "/assets/") || strings.HasPrefix(p, "/upload/"):
			assets.ServeHTTP(w, r)
		case !strings.Contains(p[1:], "/"):
			name := strings.TrimSuffix(p[1:], ".html")
			if name == "" {
				name = "index"
			}
			if _, ok := pageKeys[name]; !ok || !isFile(filepath.Join(a.cfg.SiteDir, name+".html")) {
				http.NotFound(w, r)
				return
			}
			// one URL per page: /about.html and /index.html redirect to the clean URL
			if clean := map[bool]string{true: "/", false: "/" + name}[name == "index"]; p != clean {
				if r.URL.RawQuery != "" {
					clean += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, clean, http.StatusMovedPermanently)
				return
			}
			a.servePage(w, r, name)
		default:
			http.NotFound(w, r)
		}
	})
}

// files serves regular files under root (never directory listings).
// Uploads and images get a long cache; everything else revalidates.
func (a *App) files(root, prefix string, immutable bool) http.Handler {
	fsrv := http.StripPrefix(strings.TrimSuffix(prefix, "/"), http.FileServer(http.Dir(root)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(path.Clean("/"+r.URL.Path), strings.TrimSuffix(prefix, "/"))
		if strings.Contains(rel, "/.") || !isFile(filepath.Join(root, filepath.FromSlash(rel))) {
			http.NotFound(w, r)
			return
		}
		if immutable || strings.HasPrefix(rel, "/upload/") || strings.HasPrefix(rel, "/assets/fonts/") || strings.HasPrefix(rel, "/assets/video/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fsrv.ServeHTTP(w, r)
	})
}

func headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if strings.HasPrefix(r.URL.Path, "/admin") {
			h.Set("X-Frame-Options", "DENY")
			h.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: blob:; media-src 'self' blob:; frame-src https://www.youtube-nocookie.com; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
		}
		next.ServeHTTP(w, r)
	})
}

// ---------- templates ----------

type Page struct {
	Title    string
	Active   string
	LoggedIn bool
	CSRF     string
	Flash    string
	Error    string
	Data     any
}

func parseTemplates() (map[string]*template.Template, error) {
	names, err := fs.Glob(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	pages := map[string]*template.Template{}
	for _, n := range names {
		base := path.Base(n)
		if base == "layout.html" {
			continue
		}
		t, err := template.New("").Funcs(funcs).ParseFS(templateFS, "templates/layout.html", n)
		if err != nil {
			return nil, err
		}
		pages[strings.TrimSuffix(base, ".html")] = t
	}
	return pages, nil
}

var roleNames = []struct{ Key, Label string }{
	{"coord", "Постановщик трюков"}, {"perf", "Каскадёр"}, {"horse", "Хорс-мастер"},
	{"rig", "Риггер"}, {"actor", "Актёр"}, {"actress", "Актриса"},
}

var funcs = template.FuncMap{
	"roles": func() any { return roleNames },
	"has": func(list []string, v string) bool {
		for _, x := range list {
			if x == v {
				return true
			}
		}
		return false
	},
	"roleLabel": func(k string) string {
		for _, r := range roleNames {
			if r.Key == k {
				return r.Label
			}
		}
		return k
	},
	"inc": func(i int) int { return i + 1 },
	"dict": func(kv ...any) map[string]any {
		m := map[string]any{}
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		return m
	},
	"when":  func(t time.Time) string { return t.Format("02.01.2006 15:04:05") },
	"kb":    func(n int64) int64 { return (n + 1023) / 1024 },
	"lines": func(s []string) string { return strings.Join(s, "\n") },
}

func (a *App) render(w http.ResponseWriter, r *http.Request, name string, p Page) {
	if s, ok := r.Context().Value(ctxKey{}).(session); ok {
		p.LoggedIn, p.CSRF = true, a.auth.csrf(s)
	}
	if p.Flash == "" && p.Error == "" && r.URL.Query().Get("saved") == "1" {
		p.Flash = "Сохранено. Изменения уже на сайте."
	}
	var buf bytes.Buffer
	if err := a.pages[name].ExecuteTemplate(&buf, "layout", p); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if p.Error != "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	buf.WriteTo(w)
}

func isFile(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}
