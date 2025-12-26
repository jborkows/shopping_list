package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel/attribute"

	_ "modernc.org/sqlite"

	"shopping/internal/db"
	"shopping/internal/domain/admin"
	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
)

type Repo struct {
	db *sql.DB
	q  *db.Queries
}

func Open(dsn string) (*sql.DB, error) {
	if err := ensureSQLiteDir(dsn); err != nil {
		return nil, err
	}

	conn, err := otelsql.Open(
		"sqlite",
		dsn,
		otelsql.WithAttributes(attribute.String("db.system", "sqlite")),
	)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func NewRepo(conn *sql.DB) *Repo {
	return &Repo{db: conn, q: db.New(conn)}
}

func (r *Repo) ListGroups(ctx context.Context) ([]products.Group, error) {
	rows, err := r.q.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]products.Group, 0, len(rows))
	for _, g := range rows {
		out = append(out, products.Group{ID: products.GroupID(g.ID), Name: g.Name})
	}
	return out, nil
}

func (r *Repo) ListProducts(ctx context.Context, filter products.ProductFilter) ([]products.Product, error) {
	limit := filter.Limit
	if limit <= 0 || limit > products.MaxProductsPageSize {
		limit = products.MaxProductsPageSize
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	onlyMissingOrLow := int64(0)
	if filter.OnlyMissingOrLow {
		onlyMissingOrLow = 1
	}

	name := strings.TrimSpace(filter.NameQuery)

	groupIDsCount := int64(len(filter.GroupIDs))
	groupIDs := make([]interface{}, 0, len(filter.GroupIDs))
	for _, gid := range filter.GroupIDs {
		groupIDs = append(groupIDs, int64(gid))
	}

	rows, err := r.q.ListProductsFiltered(ctx, db.ListProductsFilteredParams{
		Column1:  onlyMissingOrLow,
		Column2:  name,
		LOWER:    name,
		Column4:  groupIDsCount,
		GroupIds: groupIDs,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	out := make([]products.Product, 0, len(rows))
	for _, p := range rows {
		var gid *products.GroupID
		if v, ok := p.GroupID.(int64); ok {
			g := products.GroupID(v)
			gid = &g
		}
		groupName := ""
		if p.GroupName.Valid {
			groupName = p.GroupName.String
		}
		out = append(out, products.Product{
			ID:          products.ProductID(p.ID),
			Name:        p.Name,
			IconKey:     p.IconKey,
			GroupID:     gid,
			GroupName:   groupName,
			Quantity:    p.QuantityValue,
			Unit:        products.Unit(p.QuantityUnit),
			MinQuantity: p.MinQuantityValue,
			Missing:     p.Missing != 0,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return out, nil
}

func (r *Repo) SuggestProductsByName(ctx context.Context, query string, limit int64) ([]products.Product, error) {
	if limit <= 0 {
		limit = 8
	}
	rows, err := r.q.SuggestProductsByName(ctx, db.SuggestProductsByNameParams{
		LOWER: query,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]products.Product, 0, len(rows))
	for _, p := range rows {
		out = append(out, products.Product{
			ID:      products.ProductID(p.ID),
			Name:    p.Name,
			IconKey: p.IconKey,
			Unit:    products.Unit(p.QuantityUnit),
		})
	}
	return out, nil
}

func (r *Repo) CountProducts(ctx context.Context, filter products.ProductFilter) (int64, error) {
	onlyMissingOrLow := int64(0)
	if filter.OnlyMissingOrLow {
		onlyMissingOrLow = 1
	}

	name := strings.TrimSpace(filter.NameQuery)

	groupIDsCount := int64(len(filter.GroupIDs))
	groupIDs := make([]interface{}, 0, len(filter.GroupIDs))
	for _, gid := range filter.GroupIDs {
		groupIDs = append(groupIDs, int64(gid))
	}

	return r.q.CountProductsFiltered(ctx, db.CountProductsFilteredParams{
		Column1:  onlyMissingOrLow,
		Column2:  name,
		LOWER:    name,
		Column4:  groupIDsCount,
		GroupIds: groupIDs,
	})
}

func (r *Repo) CreateGroup(ctx context.Context, name string) (products.GroupID, error) {
	id, err := r.q.CreateGroup(ctx, name)
	return products.GroupID(id), err
}

func (r *Repo) CreateProduct(ctx context.Context, p products.NewProduct) (products.ProductID, error) {
	if p.Name == "" {
		return 0, errors.New("name required")
	}
	var gid any = nil
	if p.GroupID != nil {
		gid = int64(*p.GroupID)
	}
	id, err := r.q.CreateProduct(ctx, db.CreateProductParams{
		Name:             p.Name,
		IconKey:          p.IconKey,
		GroupID:          gid,
		QuantityValue:    p.Quantity,
		QuantityUnit:     string(p.Unit),
		MinQuantityValue: p.MinQuantity,
	})
	return products.ProductID(id), err
}

func (r *Repo) SetProductQuantity(ctx context.Context, productID products.ProductID, qty float64) error {
	return r.q.SetProductQuantity(ctx, db.SetProductQuantityParams{QuantityValue: qty, ID: int64(productID)})
}

func (r *Repo) SetProductMinQuantity(ctx context.Context, productID products.ProductID, min float64) error {
	return r.q.SetProductMinQuantity(ctx, db.SetProductMinQuantityParams{MinQuantityValue: min, ID: int64(productID)})
}

func (r *Repo) SetProductMissing(ctx context.Context, productID products.ProductID, missing bool) error {
	v := int64(0)
	if missing {
		v = 1
	}
	return r.q.SetProductMissing(ctx, db.SetProductMissingParams{Missing: v, ID: int64(productID)})
}

func (r *Repo) SetProductGroup(ctx context.Context, productID products.ProductID, groupID *products.GroupID) error {
	var gid any = nil
	if groupID != nil {
		gid = int64(*groupID)
	}
	return r.q.SetProductGroup(ctx, db.SetProductGroupParams{GroupID: gid, ID: int64(productID)})
}

func (r *Repo) SetProductUnit(ctx context.Context, productID products.ProductID, unit products.Unit) error {
	return r.q.SetProductUnit(ctx, db.SetProductUnitParams{QuantityUnit: string(unit), ID: int64(productID)})
}

func (r *Repo) ResolveIconKeyForName(ctx context.Context, name string) (string, bool, error) {
	iconKey, err := r.q.ResolveProductIconKeyByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return iconKey, true, nil
}

func (r *Repo) OptimizeDB(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "PRAGMA optimize")
	return err
}

func (r *Repo) ListItems(ctx context.Context) ([]shoppinglist.Item, error) {
	rows, err := r.q.ListShoppingListItems(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]shoppinglist.Item, 0, len(rows))
	for _, item := range rows {
		var pid *products.ProductID
		if item.ProductID.Valid {
			v := products.ProductID(item.ProductID.Int64)
			pid = &v
		}
		out = append(out, shoppinglist.Item{
			ID:        shoppinglist.ItemID(item.ID),
			Name:      item.Name,
			ProductID: pid,
			IconKey:   item.IconKey,
			GroupName: item.GroupName,
			Quantity:  item.QuantityValue,
			Unit:      products.Unit(item.QuantityUnit),
			Done:      item.Done != 0,
			CreatedAt: item.CreatedAt,
		})
	}
	return out, nil
}

func (r *Repo) GetItem(ctx context.Context, id shoppinglist.ItemID) (shoppinglist.Item, error) {
	row, err := r.q.GetShoppingListItem(ctx, int64(id))
	if err != nil {
		return shoppinglist.Item{}, err
	}
	var pid *products.ProductID
	if row.ProductID.Valid {
		v := products.ProductID(row.ProductID.Int64)
		pid = &v
	}
	return shoppinglist.Item{
		ID:        shoppinglist.ItemID(row.ID),
		Name:      row.Name,
		ProductID: pid,
		IconKey:   row.IconKey,
		GroupName: row.GroupName,
		Quantity:  row.QuantityValue,
		Unit:      products.Unit(row.QuantityUnit),
		Done:      row.Done != 0,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *Repo) AddItemByName(ctx context.Context, name string, qty float64, unit products.Unit) error {
	return r.q.AddShoppingListItemByName(ctx, db.AddShoppingListItemByNameParams{
		Name:          name,
		QuantityValue: qty,
		QuantityUnit:  string(unit),
	})
}

func (r *Repo) AddItemByProductID(ctx context.Context, productID int64) error {
	return r.q.AddShoppingListItemByProductID(ctx, productID)
}

func (r *Repo) SetDone(ctx context.Context, id shoppinglist.ItemID, done bool) error {
	v := int64(0)
	if done {
		v = 1
	}
	return r.q.SetShoppingListItemDone(ctx, db.SetShoppingListItemDoneParams{Done: v, ID: int64(id)})
}

func (r *Repo) SetQuantity(ctx context.Context, id shoppinglist.ItemID, qty float64, unit products.Unit) error {
	return r.q.SetShoppingListItemQuantity(ctx, db.SetShoppingListItemQuantityParams{QuantityValue: qty, QuantityUnit: string(unit), ID: int64(id)})
}

func (r *Repo) Delete(ctx context.Context, id shoppinglist.ItemID) error {
	return r.q.DeleteShoppingListItem(ctx, int64(id))
}

func (r *Repo) LinkToProduct(ctx context.Context, id shoppinglist.ItemID, productID int64, name string) error {
	return r.q.LinkShoppingListItemToProduct(ctx, db.LinkShoppingListItemToProductParams{
		ProductID: sql.NullInt64{Int64: productID, Valid: true},
		Name:      name,
		ID:        int64(id),
	})
}

func (r *Repo) FindProductIDByName(ctx context.Context, name string) (int64, bool, error) {
	id, err := r.q.FindProductIDByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

var _ products.Queries = (*Repo)(nil)
var _ admin.Maintenance = (*Repo)(nil)
var _ products.Repository = (*Repo)(nil)
var _ shoppinglist.Repository = (*Repo)(nil)

func ensureSQLiteDir(dsn string) error {
	path := sqlitePathFromDSN(dsn)
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func sqlitePathFromDSN(dsn string) string {
	if !strings.HasPrefix(dsn, "file:") {
		return ""
	}
	rest := strings.TrimPrefix(dsn, "file:")
	rest, _, _ = strings.Cut(rest, "?")
	rest = strings.TrimSpace(rest)
	if rest == "" || rest == ":memory:" {
		return ""
	}
	return filepath.Clean(rest)
}
