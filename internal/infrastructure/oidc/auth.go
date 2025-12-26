package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"shopping/internal/infrastructure/config"
)

type User struct {
	Subject string
	Email   string
	Name    string
	Admin   bool
}

type Authenticator interface {
	Middleware(next http.Handler) http.Handler
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
	CurrentUser(r *http.Request) (*User, bool)
}

func New(cfg config.Config) (Authenticator, error) {
	if cfg.AuthDisabled {
		return &disabledAuth{}, nil
	}
	return newOIDC(cfg)
}

type disabledAuth struct{}

func (a *disabledAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
}

func (a *disabledAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
func (a *disabledAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
func (a *disabledAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
func (a *disabledAuth) CurrentUser(r *http.Request) (*User, bool) {
	return &User{Subject: "dev", Email: "dev@example.com", Name: "Dev"}, true
}

type oidcAuth struct {
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier

	sessionsMu sync.RWMutex
	sessions   map[string]session

	stateCookieName string
	nonceCookieName string
	sessionCookie   string
}

type session struct {
	user      User
	expiresAt time.Time
}

func newOIDC(cfg config.Config) (*oidcAuth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, cfg.OIDCIssuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.OIDCClientID,
		ClientSecret: cfg.OIDCClientSecret,
		RedirectURL:  cfg.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &oidcAuth{
		oauth2Config:    oauth2Config,
		verifier:        provider.Verifier(&oidc.Config{ClientID: cfg.OIDCClientID}),
		sessions:        map[string]session{},
		stateCookieName: "oidc_state",
		nonceCookieName: "oidc_nonce",
		sessionCookie:   "shopping_session",
	}, nil
}

func (a *oidcAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r) {
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := a.CurrentUser(r); ok {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func isPublicPath(r *http.Request) bool {
	switch r.URL.Path {
	case "/healthz", "/login", "/oauth2/callback":
		return true
	}
	return strings.HasPrefix(r.URL.Path, "/static/")
}

func (a *oidcAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, _ := randToken(24)
	nonce, _ := randToken(24)

	secure := isSecureRequest(r)
	setCookie(w, a.stateCookieName, state, secure)
	setCookie(w, a.nonceCookieName, nonce, secure)

	http.Redirect(w, r, a.oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (a *oidcAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	state := r.URL.Query().Get("state")
	if state == "" || state != cookieValue(r, a.stateCookieName) {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	token, err := a.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusBadRequest)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		http.Error(w, "missing id_token", http.StatusBadRequest)
		return
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "invalid id_token", http.StatusBadRequest)
		return
	}

	nonce := cookieValue(r, a.nonceCookieName)
	if nonce == "" {
		http.Error(w, "missing nonce", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce {
		http.Error(w, "invalid nonce", http.StatusBadRequest)
		return
	}

	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "invalid claims", http.StatusBadRequest)
		return
	}

	var rawClaims map[string]any
	_ = idToken.Claims(&rawClaims)
	admin := extractAdminFromClaims(rawClaims)

	sid, _ := randToken(32)
	a.sessionsMu.Lock()
	a.sessions[sid] = session{
		user:      User{Subject: claims.Subject, Email: claims.Email, Name: claims.Name, Admin: admin},
		expiresAt: time.Now().Add(24 * time.Hour),
	}
	a.sessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     a.sessionCookie,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureRequest(r),
	})

	clearCookie(w, a.stateCookieName)
	clearCookie(w, a.nonceCookieName)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *oidcAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(a.sessionCookie); err == nil && c.Value != "" {
		a.sessionsMu.Lock()
		delete(a.sessions, c.Value)
		a.sessionsMu.Unlock()
	}
	clearCookie(w, a.sessionCookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (a *oidcAuth) CurrentUser(r *http.Request) (*User, bool) {
	c, err := r.Cookie(a.sessionCookie)
	if err != nil || c.Value == "" {
		return nil, false
	}

	a.sessionsMu.RLock()
	s, ok := a.sessions[c.Value]
	a.sessionsMu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(s.expiresAt) {
		a.sessionsMu.Lock()
		delete(a.sessions, c.Value)
		a.sessionsMu.Unlock()
		return nil, false
	}
	return &s.user, true
}

func randToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCookie(w http.ResponseWriter, name, value string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   600,
	})
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func cookieValue(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func isSecureRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func extractAdminFromClaims(claims map[string]any) bool {
	if claims == nil {
		return false
	}
	if v, ok := claims["admin"]; ok {
		return anyToBool(v)
	}
	return false
}

func anyToBool(v any) bool {
	switch vv := v.(type) {
	case bool:
		return vv
	case string:
		return strings.EqualFold(strings.TrimSpace(vv), "true")
	case float64:
		return vv != 0
	default:
		return false
	}
}
