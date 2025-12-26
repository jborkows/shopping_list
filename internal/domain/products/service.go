package products

import "context"

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
	if p.IconKey == "" {
		if iconKey, ok, err := s.repo.ResolveIconKeyForName(ctx, p.Name); err != nil {
			return 0, err
		} else if ok {
			p.IconKey = iconKey
		} else {
			p.IconKey = "cart"
		}
	}
	if p.Unit == "" {
		p.Unit = UnitPiece
	}
	if _, err := NormalizeUnit(p.Unit); err != nil {
		return 0, err
	}
	if p.Quantity < 0 {
		return 0, ErrQuantityMustBeNonNegative
	}
	if p.MinQuantity < 0 {
		return 0, ErrMinQuantityMustBeNonNegative
	}
	return s.repo.CreateProduct(ctx, p)
}

func (s *Service) SetProductQuantity(ctx context.Context, productID ProductID, qty float64) error {
	if qty < 0 {
		return ErrQuantityMustBeNonNegative
	}
	return s.repo.SetProductQuantity(ctx, productID, qty)
}

func (s *Service) SetProductMinQuantity(ctx context.Context, productID ProductID, min float64) error {
	if min < 0 {
		return ErrMinQuantityMustBeNonNegative
	}
	return s.repo.SetProductMinQuantity(ctx, productID, min)
}

func (s *Service) SetProductMissing(ctx context.Context, productID ProductID, missing bool) error {
	return s.repo.SetProductMissing(ctx, productID, missing)
}

func (s *Service) SetProductGroup(ctx context.Context, productID ProductID, groupID *GroupID) error {
	return s.repo.SetProductGroup(ctx, productID, groupID)
}

func (s *Service) SetProductUnit(ctx context.Context, productID ProductID, unit Unit) error {
	u, err := NormalizeUnit(unit)
	if err != nil {
		return err
	}
	return s.repo.SetProductUnit(ctx, productID, u)
}
