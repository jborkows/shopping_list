package web

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"

	"shopping/internal/domain/products"
)

func (s *Server) handleProductSuggestionsPartial(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		q = strings.TrimSpace(r.URL.Query().Get("name"))
	}
	if len([]rune(q)) < 2 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	suggestions, err := s.products.qry.SuggestProductsByName(ctx, q, 8)
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = renderProductSuggestions(w, suggestions)
}

func renderProductSuggestions(w http.ResponseWriter, suggestions []products.Product) error {
	if len(suggestions) == 0 {
		_, err := w.Write([]byte(""))
		return err
	}
	if _, err := w.Write([]byte(`<div class="suggestions">`)); err != nil {
		return err
	}
	for _, p := range suggestions {
		iconKey := strings.TrimSpace(p.IconKey)
		if iconKey == "" {
			iconKey = "cart"
		}
		id := strconv.FormatInt(int64(p.ID), 10)
		if _, err := w.Write([]byte(
			`<button class="suggestion" type="button" ` +
				`hx-post="/shopping-list/from-product/` + id + `" ` +
				`hx-target="#shopping-list" hx-swap="outerHTML">` +
				`<img class="icon" src="/static/icons/` + iconKey + `.svg" alt="" />` +
				templ.EscapeString(p.Name) +
				`</button>`)); err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(`</div>`))
	return err
}
