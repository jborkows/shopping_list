package products

import "context"

// Queries are read-side operations (CQRS).
type Queries interface {
	ListGroups(ctx context.Context) ([]Group, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]Product, error)
	SuggestProductsByName(ctx context.Context, query string, limit int64) ([]Product, error)
	CountProducts(ctx context.Context, filter ProductFilter) (int64, error)
}
