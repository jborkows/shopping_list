package products

import "context"

// Repository is the write-side persistence port used by Service.
// Infrastructure (e.g., SQLite) implements this interface.
type Repository interface {
	CreateGroup(ctx context.Context, name string) (GroupID, error)
	CreateProduct(ctx context.Context, p NewProduct) (ProductID, error)

	SetProductQuantity(ctx context.Context, productID ProductID, qty int) error
	SetProductMinQuantity(ctx context.Context, productID ProductID, min int) error
	SetProductMissing(ctx context.Context, productID ProductID, missing bool) error
	SetProductGroup(ctx context.Context, productID ProductID, groupID *GroupID) error
}
