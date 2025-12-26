package web

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleAutoIcon(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	initial := strings.ToUpper(string([]rune(name)[0]))

	sum := sha1.Sum([]byte(strings.ToLower(name)))
	h := hex.EncodeToString(sum[:])
	// Pleasant deterministic color.
	bg := "#" + h[:6]

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="%s" stop-opacity="0.9"/>
      <stop offset="1" stop-color="%s" stop-opacity="0.6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="10" fill="url(#g)" stroke="rgba(255,255,255,0.18)"/>
  <text x="16" y="20" text-anchor="middle" font-family="ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial"
        font-size="14" font-weight="700" fill="rgba(231,237,247,0.95)">%s</text>
</svg>`, bg, bg, initial)
}
