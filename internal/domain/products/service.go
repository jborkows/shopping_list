package products

import (
	"context"
	"errors"
)

// Service applies business rules and input validation for write operations.
// It delegates persistence to the underlying Repository.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGroup(ctx context.Context, name string) (GroupID, error) {
	normalized, err := NormalizeGroupName(name)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateGroup(ctx, normalized)
}

func (s *Service) CreateProduct(ctx context.Context, p NewProduct) (ProductID, error) {
	normalized, err := NormalizeProductName(p.Name)
	if err != nil {
		return 0, err
	}
	p.Name = normalized
	return s.repo.CreateProduct(ctx, p)
}

func (s *Service) SetProductQuantity(ctx context.Context, productID ProductID, qty int) error {
	if qty < 0 {
		return errors.New("quantity must be >= 0")
	}
	return s.repo.SetProductQuantity(ctx, productID, qty)
}

func (s *Service) SetProductMinQuantity(ctx context.Context, productID ProductID, min int) error {
	if min < 0 {
		return errors.New("min_quantity must be >= 0")
	}
	return s.repo.SetProductMinQuantity(ctx, productID, min)
}

func (s *Service) SetProductMissing(ctx context.Context, productID ProductID, missing bool) error {
	return s.repo.SetProductMissing(ctx, productID, missing)
}

func (s *Service) SetProductGroup(ctx context.Context, productID ProductID, groupID *GroupID) error {
	return s.repo.SetProductGroup(ctx, productID, groupID)
}
