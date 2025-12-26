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
	onlyMissing, nameQuery, groupIDs, page := parseProductsListQuery(r)
	perPage := products.MaxProductsPageSize

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	total, err := s.products.qry.CountProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
	})
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	totalPages := int64(1)
	if total > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * perPage
	list, err := s.products.qry.ListProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
		Limit:            perPage,
		Offset:           offset,
	})
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	data := productsPageData{
		Base: baseData{
			Title:   "Zapasy",
			User:    user,
			HTMXSrc: s.cfg.HTMXSrc,
			IsAdmin: s.isAdmin(user),
		},
		Groups:           groups,
		Products:         list,
		OnlyMissing:      onlyMissing,
		NameQuery:        nameQuery,
		SelectedGroupIDs: groupIDs,
		Page:             page,
		TotalPages:       totalPages,
		Total:            total,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ProductsPage(data).Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("render: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleProductsPartial(w http.ResponseWriter, r *http.Request) {
	onlyMissing, nameQuery, groupIDs, page := parseProductsListQuery(r)
	perPage := products.MaxProductsPageSize
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	groups, err := s.products.qry.ListGroups(ctx)
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	total, err := s.products.qry.CountProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
	})
	if err != nil {
		s.writeDBError(w, err)
		return
	}
	totalPages := int64(1)
	if total > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * perPage
	list, err := s.products.qry.ListProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: onlyMissing,
		NameQuery:        nameQuery,
		GroupIDs:         groupIDs,
		Limit:            perPage,
		Offset:           offset,
	})
	if err != nil {
		s.writeDBError(w, err)
		return
	}

	data := productsListData{
		Groups:           groups,
		Products:         list,
		OnlyMissing:      onlyMissing,
		NameQuery:        nameQuery,
		SelectedGroupIDs: groupIDs,
		Page:             page,
		TotalPages:       totalPages,
		Total:            total,
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
	qty, err := parseFloat(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	minQty, err := parseFloat(r.FormValue("min_quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa minimalna ilość.", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var gid *products.GroupID
	if ok {
		gid = &groupID
	}
	if _, err := s.products.svc.CreateProduct(ctx, products.NewProduct{
		Name:        name,
		GroupID:     gid,
		Quantity:    qty,
		MinQuantity: minQty,
		Unit:        unit,
	}); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.handleProductsPartial(w, r)
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
	qty, err := parseFloat(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa ilość.", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductQuantity(ctx, id, qty); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.handleProductsPartial(w, r)
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
	min, err := parseFloat(r.FormValue("min_quantity"))
	if err != nil {
		http.Error(w, "Nieprawidłowa minimalna ilość.", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductMinQuantity(ctx, id, min); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.handleProductsPartial(w, r)
}

func (s *Server) handleSetUnit(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathProductID(r, "id")
	if !ok {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	unit := products.Unit(strings.TrimSpace(r.FormValue("unit")))
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := s.products.svc.SetProductUnit(ctx, id, unit); err != nil {
		s.writeUserError(w, err)
		return
	}
	s.handleProductsPartial(w, r)
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
	s.handleProductsPartial(w, r)
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
	s.handleProductsPartial(w, r)
}

func parseFloat(v string) (float64, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}
	return strconv.ParseFloat(v, 64)
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
	if r.URL.RawQuery != "" {
		redirect += "?" + r.URL.RawQuery
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

func parseProductsListQuery(r *http.Request) (onlyMissing bool, nameQuery string, groupIDs []products.GroupID, page int64) {
	q := r.URL.Query()
	onlyMissing = q.Get("missing") == "1"
	nameQuery = strings.TrimSpace(q.Get("q"))

	seen := make(map[products.GroupID]struct{})
	for _, raw := range q["group_id"] {
		id, ok := parseOptionalGroupID(raw)
		if !ok {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		groupIDs = append(groupIDs, id)
	}

	page = 1
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		if p, err := strconv.ParseInt(raw, 10, 64); err == nil && p > 0 {
			page = p
		}
	}
	return onlyMissing, nameQuery, groupIDs, page
}
