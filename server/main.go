// Command server serves the Nomad Stunts landing page and its admin panel (/admin).
//
// Content (team, projects, gallery, showreel, contacts) lives in DATA_DIR/content.json and is
// served to the site as /assets/data.js; uploads go to DATA_DIR/media. See server/README.md.
package main

import (
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Port, SiteDir, DataDir string
	SiteURL                string // public origin for canonical/sitemap links, e.g. https://nomadstunts.com
	AdminUser, AdminPass   string
	Secret                 []byte
	TrustProxy             bool
}

func main() {
	cfg := Config{
		Port:       env("PORT", "8080"),
		SiteDir:    env("SITE_DIR", "."),
		DataDir:    env("DATA_DIR", "./data"),
		SiteURL:    strings.TrimRight(os.Getenv("SITE_URL"), "/"),
		AdminUser:  env("ADMIN_USER", "admin"),
		AdminPass:  os.Getenv("ADMIN_PASSWORD"),
		TrustProxy: os.Getenv("TRUST_PROXY") == "1",
	}
	if len(cfg.AdminPass) < 10 {
		log.Fatal("ADMIN_PASSWORD must be set (at least 10 characters)")
	}
	if s := os.Getenv("SESSION_SECRET"); s != "" {
		sum := sha256.Sum256([]byte(s))
		cfg.Secret = sum[:]
	} else {
		cfg.Secret = []byte(randHex(32))
		log.Print("SESSION_SECRET not set: admins will be logged out on every restart")
	}

	app, err := NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	log.Printf("Nomad Stunts on http://localhost:%s  (admin: /admin, data: %s)", cfg.Port, cfg.DataDir)
	log.Fatal(srv.ListenAndServe())
}

func NewApp(cfg Config) (*App, error) {
	store, err := OpenStore(cfg.DataDir, filepath.Join(cfg.SiteDir, "assets", "data.js"))
	if err != nil {
		return nil, err
	}
	media, err := NewMedia(filepath.Join(cfg.DataDir, "media"))
	if err != nil {
		return nil, err
	}
	pages, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	i18n, err := loadI18N(cfg.SiteDir)
	if err != nil {
		return nil, err
	}
	return &App{
		cfg:   cfg,
		store: store,
		media: media,
		auth:  NewAuth(cfg.AdminUser, cfg.AdminPass, cfg.Secret, cfg.TrustProxy),
		pages: pages,
		i18n:  i18n,
	}, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
