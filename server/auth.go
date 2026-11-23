package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	sessionCookie = "ns_admin"
	sessionTTL    = 7 * 24 * time.Hour
	loginWindow   = 15 * time.Minute
	loginMaxFails = 10
)

// Auth is a single admin account (from env) with stateless HMAC-signed session cookies.
type Auth struct {
	user, pass string
	secret     []byte
	trustProxy bool

	mu    sync.Mutex
	fails map[string][]time.Time
}

type session struct{ nonce string }

func NewAuth(user, pass string, secret []byte, trustProxy bool) *Auth {
	return &Auth{user: user, pass: pass, secret: secret, trustProxy: trustProxy, fails: map[string][]time.Time{}}
}

// check compares in constant time; hashing first makes the comparison length-independent.
func (a *Auth) check(user, pass string) bool {
	hu, wu := sha256.Sum256([]byte(user)), sha256.Sum256([]byte(a.user))
	hp, wp := sha256.Sum256([]byte(pass)), sha256.Sum256([]byte(a.pass))
	return subtle.ConstantTimeCompare(hu[:], wu[:])&subtle.ConstantTimeCompare(hp[:], wp[:]) == 1
}

// allowLogin reports whether this client is under the failed-login limit.
func (a *Auth) allowLogin(r *http.Request) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	ip := a.clientIP(r)
	recent := a.fails[ip][:0]
	for _, t := range a.fails[ip] {
		if time.Since(t) < loginWindow {
			recent = append(recent, t)
		}
	}
	if len(recent) == 0 {
		delete(a.fails, ip)
	} else {
		a.fails[ip] = recent
	}
	return len(recent) < loginMaxFails
}

func (a *Auth) recordFail(r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ip := a.clientIP(r)
	a.fails[ip] = append(a.fails[ip], time.Now())
}

func (a *Auth) clientIP(r *http.Request) string {
	if a.trustProxy {
		if f := r.Header.Get("X-Forwarded-For"); f != "" {
			return strings.TrimSpace(strings.Split(f, ",")[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (a *Auth) sign(payload string) string {
	m := hmac.New(sha256.New, a.secret)
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func (a *Auth) login(w http.ResponseWriter, r *http.Request) {
	nonce := randHex(16)
	payload := base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().Add(sessionTTL).Unix(), 10) + ":" + nonce))
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: payload + "." + a.sign(payload), Path: "/admin",
		MaxAge: int(sessionTTL.Seconds()), HttpOnly: true, Secure: isHTTPS(r), SameSite: http.SameSiteStrictMode,
	})
}

func (a *Auth) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/admin", MaxAge: -1, HttpOnly: true, Secure: isHTTPS(r), SameSite: http.SameSiteStrictMode})
}

func (a *Auth) session(r *http.Request) (session, bool) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return session{}, false
	}
	payload, sig, ok := strings.Cut(c.Value, ".")
	if !ok || !hmac.Equal([]byte(sig), []byte(a.sign(payload))) {
		return session{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return session{}, false
	}
	exp, nonce, ok := strings.Cut(string(raw), ":")
	unix, err := strconv.ParseInt(exp, 10, 64)
	if !ok || err != nil || time.Now().Unix() > unix {
		return session{}, false
	}
	return session{nonce: nonce}, true
}

// csrf is derived from the session, so it needs no server-side storage.
func (a *Auth) csrf(s session) string { return a.sign("csrf:" + s.nonce) }

func (a *Auth) validCSRF(r *http.Request, s session) bool {
	return hmac.Equal([]byte(r.FormValue("_csrf")), []byte(a.csrf(s)))
}

func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
