package web

import (
	"net/http"
	"strings"
)

func (s *Server) writeDBError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "no such table") {
		http.Error(w, "database schema missing: run migrations in ./migrations (see README.md)", http.StatusInternalServerError)
		return
	}
	http.Error(w, msg, http.StatusInternalServerError)
}
