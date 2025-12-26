package shoppinglist

import (
	"context"

	"shopping/internal/domain/products"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListItems(ctx context.Context) ([]Item, error) {
	return s.repo.ListItems(ctx)
}

func (s *Service) GetItem(ctx context.Context, id ItemID) (Item, error) {
	return s.repo.GetItem(ctx, id)
}

func (s *Service) AddItemByName(ctx context.Context, name string, qty float64, unit products.Unit) error {
	normalized, err := NormalizeItemName(name)
	if err != nil {
		return err
	}
	if qty <= 0 {
		return ErrQuantityMustBePositive
	}
	if unit == "" {
		unit = products.UnitPiece
	}
	if _, err := products.NormalizeUnit(unit); err != nil {
		return err
	}
	return s.repo.AddItemByName(ctx, normalized, qty, unit)
}

func (s *Service) AddItemByProductID(ctx context.Context, productID int64) error {
	return s.repo.AddItemByProductID(ctx, productID)
}

func (s *Service) SetDone(ctx context.Context, id ItemID, done bool) error {
	return s.repo.SetDone(ctx, id, done)
}

func (s *Service) SetQuantity(ctx context.Context, id ItemID, qty float64, unit products.Unit) error {
	if qty <= 0 {
		return ErrQuantityMustBePositive
	}
	if _, err := products.NormalizeUnit(unit); err != nil {
		return err
	}
	return s.repo.SetQuantity(ctx, id, qty, unit)
}

func (s *Service) Delete(ctx context.Context, id ItemID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) LinkToProduct(ctx context.Context, id ItemID, productID int64, name string) error {
	return s.repo.LinkToProduct(ctx, id, productID, name)
}

func (s *Service) FindProductIDByName(ctx context.Context, name string) (int64, bool, error) {
	return s.repo.FindProductIDByName(ctx, name)
}
