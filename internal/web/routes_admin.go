package web

import "net/http"

func (s *Server) registerAdminRoutes(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	mux.Handle("GET /admin", wrap(http.HandlerFunc(s.handleAdminPage)))
	mux.Handle("POST /admin/db/optimize", wrap(http.HandlerFunc(s.handleDBOptimize)))
}
