package shoppinglist

import (
	"context"

	"shopping/internal/domain/products"
)

type Repository interface {
	ListItems(ctx context.Context) ([]Item, error)
	GetItem(ctx context.Context, id ItemID) (Item, error)

	AddItemByName(ctx context.Context, name string, qty float64, unit products.Unit) error
	AddItemByProductID(ctx context.Context, productID int64) error

	SetDone(ctx context.Context, id ItemID, done bool) error
	SetQuantity(ctx context.Context, id ItemID, qty float64, unit products.Unit) error
	Delete(ctx context.Context, id ItemID) error

	LinkToProduct(ctx context.Context, id ItemID, productID int64, name string) error
	FindProductIDByName(ctx context.Context, name string) (int64, bool, error)
}
