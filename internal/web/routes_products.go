package web

import "net/http"

func (s *Server) registerProductRoutes(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	mux.Handle("GET /products", wrap(http.HandlerFunc(s.handleProductsPage)))
	mux.Handle("GET /partials/products", wrap(http.HandlerFunc(s.handleProductsPartial)))
	mux.Handle("POST /products", wrap(http.HandlerFunc(s.handleCreateProduct)))
	mux.Handle("POST /products/{id}/qty", wrap(http.HandlerFunc(s.handleSetQuantity)))
	mux.Handle("POST /products/{id}/min", wrap(http.HandlerFunc(s.handleSetMin)))
	mux.Handle("POST /products/{id}/missing", wrap(http.HandlerFunc(s.handleSetMissing)))
	mux.Handle("POST /products/{id}/group", wrap(http.HandlerFunc(s.handleSetGroup)))
	mux.Handle("POST /groups", wrap(http.HandlerFunc(s.handleCreateGroup)))
}
