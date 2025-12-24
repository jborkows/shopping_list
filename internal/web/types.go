package web

import (
	"shopping/internal/domain/products"
	"shopping/internal/infrastructure/oidc"
)

type baseData struct {
	Title   string
	User    *oidc.User
	HTMXSrc string
	IsAdmin bool
}

type productsPageData struct {
	Base        baseData
	Groups      []products.Group
	Products    []products.Product
	OnlyMissing bool
}

type productsListData struct {
	Groups      []products.Group
	Products    []products.Product
	OnlyMissing bool
}
