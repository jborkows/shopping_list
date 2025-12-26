package web

import (
	"net/http"

	"shopping/internal/domain/admin"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/config"
	"shopping/internal/infrastructure/oidc"
)

type Server struct {
	cfg      config.Config
	auth     oidc.Authenticator
	products productsComponent
	shopping shoppingComponent
	admin    adminComponent
}

func NewServer(cfg config.Config, qry products.Queries, svc *products.Service, shoppingSvc *shoppinglist.Service, adminMaintenance admin.Maintenance, authenticator oidc.Authenticator) *Server {
	return &Server{
		cfg:  cfg,
		auth: authenticator,
		products: productsComponent{
			qry: qry,
			svc: svc,
		},
		shopping: shoppingComponent{
			svc: shoppingSvc,
		},
		admin: adminComponent{
			maintenance: adminMaintenance,
		},
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.Handle("GET /login", http.HandlerFunc(s.auth.HandleLogin))
	mux.Handle("GET /oauth2/callback", http.HandlerFunc(s.auth.HandleCallback))
	mux.Handle("POST /logout", http.HandlerFunc(s.auth.HandleLogout))

	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/products", http.StatusFound) }))

	wrap := s.auth.Middleware

	s.registerProductRoutes(mux, wrap)
	s.registerAdminRoutes(mux, wrap)

	return mux
}
