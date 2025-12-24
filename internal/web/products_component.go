package web

import "shopping/internal/domain/products"

type productsComponent struct {
	qry products.Queries
	svc *products.Service
}
