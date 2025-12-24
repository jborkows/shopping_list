package web

import "net/http"

func (s *Server) registerAdminRoutes(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	mux.Handle("POST /admin/db/optimize", wrap(http.HandlerFunc(s.handleDBOptimize)))
}
