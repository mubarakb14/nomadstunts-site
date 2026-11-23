package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Content is everything the landing page renders from window.NS_DATA.
// JSON keys match what assets/app.js reads, so the site needs no translation layer.
type Content struct {
	Showreel Showreel `json:"showreel"`
	Contacts Contacts `json:"contacts"`
	Founder  Founder  `json:"founder"`
	Team     []Person `json:"team"`
	Films    []Film   `json:"films"`
	Gallery  []Photo  `json:"gallery"`
}

type L10n struct {
	En string `json:"en,omitempty"`
	Ru string `json:"ru,omitempty"`
	Kk string `json:"kk,omitempty"`
}

// Name holds [first, last] per language.
type Name struct {
	En [2]string `json:"en"`
	Ru [2]string `json:"ru"`
	Kk [2]string `json:"kk"`
}

type Showreel struct {
	MP4     string `json:"mp4,omitempty"`
	Poster  string `json:"poster,omitempty"`
	YouTube string `json:"youtube,omitempty"`
}

type Contacts struct {
	Phones []string `json:"phones"`
	Email  string   `json:"email"`
	Base   L10n     `json:"base"`
	Map    string   `json:"map"`
}

type Founder struct {
	Name  Name   `json:"n"`
	Photo string `json:"photo"`
	Role  L10n   `json:"role"`
	Hon   L10n   `json:"h"`
}

type Person struct {
	ID    string   `json:"id"`
	Name  Name     `json:"n"`
	Img   string   `json:"img"`
	Photo string   `json:"photo"`
	Roles []string `json:"roles"`
	Hon   L10n     `json:"h"`
}

type FilmInfo struct {
	Type     string `json:"type,omitempty"`
	Country  string `json:"country,omitempty"`
	Director string `json:"director,omitempty"`
	Stars    string `json:"stars,omitempty"`
}

type Film struct {
	ID       string   `json:"id"`
	Year     int      `json:"y"`
	Poster   string   `json:"poster"`
	Img      string   `json:"img"`
	Title    L10n     `json:"t"`
	En       FilmInfo `json:"en"`
	Ru       FilmInfo `json:"ru"`
	Kk       FilmInfo `json:"kk"`
	Featured int      `json:"featured,omitempty"` // position in the home-page reel; 0 = not pinned
}

type Photo struct {
	ID    string `json:"id"`
	Img   string `json:"img"`
	Thumb string `json:"thumb"`
}

const historyKeep = 60

var snapshotName = regexp.MustCompile(`^content-\d{8}-\d{6}\.\d{3}\.json$`)

// Store keeps Content in memory and in DATA_DIR/content.json.
// Every save first copies the previous file into history/, so any edit can be undone.
type Store struct {
	mu      sync.RWMutex
	dir     string
	content Content
	js      []byte
	updated time.Time
}

// OpenStore loads DATA_DIR/content.json, seeding it from the site's data.js on first run.
func OpenStore(dir, seedJS string) (*Store, error) {
	if err := os.MkdirAll(filepath.Join(dir, "history"), 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir}
	raw, err := os.ReadFile(s.path())
	if errors.Is(err, os.ErrNotExist) {
		if raw, err = readDataJS(seedJS); err != nil {
			return nil, fmt.Errorf("seed from %s: %w", seedJS, err)
		}
		var c Content
		if err := decodeStrict(raw, &c); err != nil {
			return nil, fmt.Errorf("seed from %s: %w", seedJS, err)
		}
		if err := s.write(c); err != nil {
			return nil, err
		}
		s.set(c)
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var c Content
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("%s: %w", s.path(), err)
	}
	s.set(c)
	if st, err := os.Stat(s.path()); err == nil {
		s.updated = st.ModTime()
	}
	return s, nil
}

func (s *Store) path() string { return filepath.Join(s.dir, "content.json") }

// Get returns a deep copy, so callers can't mutate the live content.
func (s *Store) Get() Content {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return clone(s.content)
}

// Modified is when the content last changed (for sitemap lastmod).
func (s *Store) Modified() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updated
}

// DataJS is the rendered assets/data.js for the public site.
func (s *Store) DataJS() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.js
}

// Update applies fn to a copy of the content and persists it. If fn fails, nothing changes.
func (s *Store) Update(fn func(*Content) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := clone(s.content)
	if err := fn(&c); err != nil {
		return err
	}
	if err := s.backup(); err != nil {
		return fmt.Errorf("backup: %w", err)
	}
	if err := s.write(c); err != nil {
		return err
	}
	s.set(c)
	return nil
}

type Snapshot struct {
	Name string
	Time time.Time
	Size int64
}

// History lists saved versions, newest first.
func (s *Store) History() ([]Snapshot, error) {
	entries, err := os.ReadDir(filepath.Join(s.dir, "history"))
	if err != nil {
		return nil, err
	}
	var out []Snapshot
	for _, e := range entries {
		if !snapshotName.MatchString(e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		t, _ := time.ParseInLocation("20060102-150405.000", strings.TrimSuffix(strings.TrimPrefix(e.Name(), "content-"), ".json"), time.Local)
		out = append(out, Snapshot{Name: e.Name(), Time: t, Size: info.Size()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name > out[j].Name })
	return out, nil
}

// Restore replaces the content with a saved version (the current one is backed up first).
func (s *Store) Restore(name string) error {
	if !snapshotName.MatchString(name) {
		return errors.New("неизвестная версия")
	}
	raw, err := os.ReadFile(filepath.Join(s.dir, "history", name))
	if err != nil {
		return err
	}
	var old Content
	if err := json.Unmarshal(raw, &old); err != nil {
		return err
	}
	return s.Update(func(c *Content) error { *c = old; return nil })
}

func (s *Store) set(c Content) {
	normalize(&c)
	s.content = c
	s.updated = time.Now()
	b, _ := json.Marshal(c) // Marshal escapes <, >, & so the payload is safe inside a script
	s.js = append(append([]byte("window.NS_DATA = "), b...), ";\n"...)
}

// normalize turns nil slices into empty ones: the site calls .map/.filter on them, and null would break it.
func normalize(c *Content) {
	if c.Team == nil {
		c.Team = []Person{}
	}
	if c.Films == nil {
		c.Films = []Film{}
	}
	if c.Gallery == nil {
		c.Gallery = []Photo{}
	}
	if c.Contacts.Phones == nil {
		c.Contacts.Phones = []string{}
	}
	for i := range c.Team {
		if c.Team[i].Roles == nil {
			c.Team[i].Roles = []string{}
		}
	}
}

func (s *Store) write(c Content) error {
	normalize(&c)
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(s.path(), b)
}

func (s *Store) backup() error {
	cur, err := os.ReadFile(s.path())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	name := "content-" + time.Now().Format("20060102-150405.000") + ".json"
	if err := writeFileAtomic(filepath.Join(s.dir, "history", name), cur); err != nil {
		return err
	}
	snaps, err := s.History()
	if err != nil {
		return err
	}
	for _, old := range snaps[min(len(snaps), historyKeep):] {
		os.Remove(filepath.Join(s.dir, "history", old.Name))
	}
	return nil
}

// writeFileAtomic writes via a temp file + rename, so a crash never leaves a half-written file.
func writeFileAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// readDataJS extracts the JSON object from "window.NS_DATA = {...};".
func readDataJS(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i, j := bytes.IndexByte(b, '{'), bytes.LastIndexByte(b, '}')
	if i < 0 || j < i {
		return nil, errors.New("no JSON object found")
	}
	return b[i : j+1], nil
}

// decodeStrict fails on unknown fields, so seeding never silently drops site data.
func decodeStrict(b []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(v)
}

func clone(c Content) Content {
	b, _ := json.Marshal(c)
	var out Content
	_ = json.Unmarshal(b, &out)
	return out
}
