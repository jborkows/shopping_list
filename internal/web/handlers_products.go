package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"shopping/internal/domain/products"
)

func (s *Server) handleProductsPage(w http.ResponseWriter, r *http.Request) {
	user, _ := s.auth.CurrentUser(r)
	onlyMissing := r.URL.Query().Get("missing") == "1"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	list, err := s.products.qry.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: onlyMissing})
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	data := productsPageData{
		Base: baseData{
			Title:   "Products",
			User:    user,
			HTMXSrc: s.cfg.HTMXSrc,
			IsAdmin: s.isAdmin(user),
		},
		Groups:      groups,
		Products:    list,
		OnlyMissing: onlyMissing,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ProductsPage(data).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleProductsPartial(w http.ResponseWriter, r *http.Request) {
	onlyMissing := r.URL.Query().Get("missing") == "1"
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	list, err := s.products.qry.ListProducts(ctx, products.ProductFilter{OnlyMissingOrLow: onlyMissing})
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	data := productsListData{
		Groups:      groups,
		Products:    list,
		OnlyMissing: onlyMissing,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ProductsList(data).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	groupID, ok := parseOptionalGroupID(r.FormValue("group_id"))

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var gid *products.GroupID
	if ok {
		gid = &groupID
	}
	if _, err := s.products.svc.CreateProduct(ctx, products.NewProduct{Name: name, GroupID: gid}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.handleProductsPartial(w, withQuery(r, "missing", r.URL.Query().Get("missing")))
}

func (s *Server) handleSetQuantity(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	qty, err := strconv.Atoi(strings.TrimSpace(r.FormValue("quantity")))
	if err != nil {
		http.Error(w, "bad quantity", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductQuantity(ctx, id, qty); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.handleProductsPartial(w, withQuery(r, "missing", r.URL.Query().Get("missing")))
}

func (s *Server) handleSetMin(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	min, err := strconv.Atoi(strings.TrimSpace(r.FormValue("min_quantity")))
	if err != nil {
		http.Error(w, "bad min_quantity", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductMinQuantity(ctx, id, min); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.handleProductsPartial(w, withQuery(r, "missing", r.URL.Query().Get("missing")))
}

func (s *Server) handleSetMissing(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	missing := r.FormValue("missing") == "on" || r.FormValue("missing") == "1" || r.FormValue("missing") == "true"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductMissing(ctx, id, missing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.handleProductsPartial(w, withQuery(r, "missing", r.URL.Query().Get("missing")))
}

func (s *Server) handleSetGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	groupID, ok := parseOptionalGroupID(r.FormValue("group_id"))
	var gid *products.GroupID
	if ok {
		gid = &groupID
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductGroup(ctx, id, gid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.handleProductsPartial(w, withQuery(r, "missing", r.URL.Query().Get("missing")))
}

func (s *Server) handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if _, err := s.products.svc.CreateGroup(ctx, name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	redirect := "/products"
	if r.URL.Query().Get("missing") == "1" {
		redirect = "/products?missing=1"
	}
	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(http.StatusNoContent)
}

func parsePathProductID(r *http.Request, param string) (products.ProductID, bool) {
	v := r.PathValue(param)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	return products.ProductID(id), err == nil
}

func parseOptionalGroupID(v string) (products.GroupID, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(v, 10, 64)
	return products.GroupID(id), err == nil
}

func withQuery(r *http.Request, key, value string) *http.Request {
	if value == "" {
		return r
	}
	r2 := r.Clone(r.Context())
	q := r2.URL.Query()
	q.Set(key, value)
	r2.URL.RawQuery = q.Encode()
	return r2
}
