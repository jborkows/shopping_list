package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Addr         string
	BaseURL      string
	DBDSN        string
	AuthDisabled bool
	AdminEmails  []string

	OIDCIssuer       string
	OIDCClientID     string
	OIDCClientSecret string
	OIDCRedirectURL  string

	HTMXSrc string
}

func FromEnv() (Config, error) {
	_ = loadDotEnv(".env")

	cfg := Config{
		Addr:    envOr("ADDR", ":8080"),
		BaseURL: envOr("BASE_URL", "http://localhost:8080"),
		DBDSN:   envOr("DB_DSN", "file:data/shopping.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)"),
		HTMXSrc: envOr("HTMX_SRC", "https://unpkg.com/htmx.org@1.9.12"),
	}

	cfg.AuthDisabled = os.Getenv("AUTH_DISABLED") == "1" || os.Getenv("AUTH_DISABLED") == "true"
	cfg.AdminEmails = splitCSV(os.Getenv("ADMIN_EMAILS"))

	cfg.OIDCIssuer = os.Getenv("OIDC_ISSUER")
	cfg.OIDCClientID = os.Getenv("OIDC_CLIENT_ID")
	cfg.OIDCClientSecret = os.Getenv("OIDC_CLIENT_SECRET")
	cfg.OIDCRedirectURL = envOr("OIDC_REDIRECT_URL", cfg.BaseURL+"/oauth2/callback")

	if !cfg.AuthDisabled {
		if cfg.OIDCIssuer == "" || cfg.OIDCClientID == "" || cfg.OIDCClientSecret == "" || cfg.OIDCRedirectURL == "" {
			return Config{}, errors.New("OIDC is enabled but required env vars are missing (OIDC_ISSUER, OIDC_CLIENT_ID, OIDC_CLIENT_SECRET, OIDC_REDIRECT_URL)")
		}
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if key == "" {
			continue
		}
		if os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, val)
	}
	return scanner.Err()
}

func splitCSV(s string) []string {
	var out []string
	for part := range strings.SplitSeq(s, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}
