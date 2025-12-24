package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"shopping/internal/infrastructure/oidc"
)

func (s *Server) isAdmin(user *oidc.User) bool {
	if user == nil {
		return false
	}
	if s.cfg.AuthDisabled {
		return true
	}
	if len(s.cfg.AdminEmails) == 0 {
		return false
	}
	for _, email := range s.cfg.AdminEmails {
		if strings.EqualFold(email, user.Email) {
			return true
		}
	}
	return false
}

func (s *Server) handleDBOptimize(w http.ResponseWriter, r *http.Request) {
	user, _ := s.auth.CurrentUser(r)
	if !s.isAdmin(user) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err := s.admin.maintenance.OptimizeDB(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("Optimized."))
}
