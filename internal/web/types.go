package web

import (
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
	"shopping/internal/infrastructure/oidc"
)

type baseData struct {
	Title   string
	User    *oidc.User
	HTMXSrc string
	IsAdmin bool
}

type productsPageData struct {
	Base             baseData
	Groups           []products.Group
	Products         []products.Product
	OnlyMissing      bool
	NameQuery        string
	SelectedGroupIDs []products.GroupID
	Page             int64
	TotalPages       int64
	Total            int64
}

type productsListData struct {
	Groups           []products.Group
	Products         []products.Product
	OnlyMissing      bool
	NameQuery        string
	SelectedGroupIDs []products.GroupID
	Page             int64
	TotalPages       int64
	Total            int64
}

type adminPageData struct {
	Base baseData
}

type shoppingListData struct {
	Items []shoppinglist.Item
}

type shoppingListPageData struct {
	Base  baseData
	Items []shoppinglist.Item
}
