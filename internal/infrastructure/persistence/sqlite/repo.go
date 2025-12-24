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
	if filter.OnlyMissingOrLow {
		rows, err := r.q.ListProductsMissingOrLow(ctx)
		if err != nil {
			return nil, err
		}
		return mapProductsMissing(rows), nil
	}

	rows, err := r.q.ListProductsAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapProductsAll(rows), nil
}

func (r *Repo) CreateGroup(ctx context.Context, name string) (products.GroupID, error) {
	id, err := r.q.CreateGroup(ctx, name)
	return products.GroupID(id), err
}

func (r *Repo) CreateProduct(ctx context.Context, p products.NewProduct) (products.ProductID, error) {
	if p.Name == "" {
		return 0, errors.New("name required")
	}
	arg := db.CreateProductParams{Name: p.Name}
	if p.GroupID != nil {
		arg.GroupID = sql.NullInt64{Int64: int64(*p.GroupID), Valid: true}
	}
	id, err := r.q.CreateProduct(ctx, arg)
	return products.ProductID(id), err
}

func (r *Repo) SetProductQuantity(ctx context.Context, productID products.ProductID, qty int) error {
	return r.q.SetProductQuantity(ctx, db.SetProductQuantityParams{ID: int64(productID), Quantity: int64(qty)})
}

func (r *Repo) SetProductMinQuantity(ctx context.Context, productID products.ProductID, min int) error {
	return r.q.SetProductMinQuantity(ctx, db.SetProductMinQuantityParams{ID: int64(productID), MinQuantity: int64(min)})
}

func (r *Repo) SetProductMissing(ctx context.Context, productID products.ProductID, missing bool) error {
	v := int64(0)
	if missing {
		v = 1
	}
	return r.q.SetProductMissing(ctx, db.SetProductMissingParams{ID: int64(productID), Missing: v})
}

func (r *Repo) SetProductGroup(ctx context.Context, productID products.ProductID, groupID *products.GroupID) error {
	arg := db.SetProductGroupParams{ID: int64(productID)}
	if groupID != nil {
		arg.GroupID = sql.NullInt64{Int64: int64(*groupID), Valid: true}
	}
	return r.q.SetProductGroup(ctx, arg)
}

func (r *Repo) OptimizeDB(ctx context.Context) error {
	return r.q.PragmaOptimize(ctx)
}

func mapProductsAll(rows []db.ListProductsAllRow) []products.Product {
	out := make([]products.Product, 0, len(rows))
	for _, p := range rows {
		out = append(out, mapProductRow(p.ID, p.Name, p.GroupID, p.GroupName, p.Quantity, p.MinQuantity, p.Missing, p.UpdatedAt))
	}
	return out
}

func mapProductsMissing(rows []db.ListProductsMissingOrLowRow) []products.Product {
	out := make([]products.Product, 0, len(rows))
	for _, p := range rows {
		out = append(out, mapProductRow(p.ID, p.Name, p.GroupID, p.GroupName, p.Quantity, p.MinQuantity, p.Missing, p.UpdatedAt))
	}
	return out
}

func mapProductRow(id int64, name string, gid sql.NullInt64, gname sql.NullString, qty, min, missing int64, updated time.Time) products.Product {
	var groupID *products.GroupID
	var groupName string
	if gid.Valid {
		v := products.GroupID(gid.Int64)
		groupID = &v
	}
	if gname.Valid {
		groupName = gname.String
	}

	return products.Product{
		ID:          products.ProductID(id),
		Name:        name,
		GroupID:     groupID,
		GroupName:   groupName,
		Quantity:    int(qty),
		MinQuantity: int(min),
		Missing:     missing == 1,
		UpdatedAt:   updated,
	}
}

var _ products.Queries = (*Repo)(nil)
var _ admin.Maintenance = (*Repo)(nil)
var _ products.Repository = (*Repo)(nil)

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
