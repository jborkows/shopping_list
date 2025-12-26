package web

import "net/http"

func (s *Server) registerProductRoutes(mux *http.ServeMux, wrap func(http.Handler) http.Handler) {
	mux.Handle("GET /products", wrap(http.HandlerFunc(s.handleProductsPage)))
	mux.Handle("GET /shopping-list", wrap(http.HandlerFunc(s.handleShoppingListPage)))
	mux.Handle("GET /icons/auto", wrap(http.HandlerFunc(s.handleAutoIcon)))
	mux.Handle("GET /partials/products", wrap(http.HandlerFunc(s.handleProductsPartial)))
	mux.Handle("GET /partials/shopping-list", wrap(http.HandlerFunc(s.handleShoppingListPartial)))
	mux.Handle("GET /partials/product-suggestions", wrap(http.HandlerFunc(s.handleProductSuggestionsPartial)))
	mux.Handle("POST /products", wrap(http.HandlerFunc(s.handleCreateProduct)))
	mux.Handle("POST /products/{id}/qty", wrap(http.HandlerFunc(s.handleSetQuantity)))
	mux.Handle("POST /products/{id}/min", wrap(http.HandlerFunc(s.handleSetMin)))
	mux.Handle("POST /products/{id}/unit", wrap(http.HandlerFunc(s.handleSetUnit)))
	mux.Handle("POST /products/{id}/missing", wrap(http.HandlerFunc(s.handleSetMissing)))
	mux.Handle("POST /products/{id}/group", wrap(http.HandlerFunc(s.handleSetGroup)))
	mux.Handle("POST /groups", wrap(http.HandlerFunc(s.handleCreateGroup)))

	mux.Handle("POST /shopping-list", wrap(http.HandlerFunc(s.handleAddShoppingListByName)))
	mux.Handle("POST /shopping-list/from-product/{id}", wrap(http.HandlerFunc(s.handleAddShoppingListFromProduct)))
	mux.Handle("PATCH /shopping-list/items/{id}", wrap(http.HandlerFunc(s.handleSetShoppingListDone)))
	mux.Handle("DELETE /shopping-list/items/{id}", wrap(http.HandlerFunc(s.handleDeleteShoppingListItem)))
	mux.Handle("POST /shopping-list/items/{id}/product", wrap(http.HandlerFunc(s.handleAddShoppingItemToSupplies)))
}
