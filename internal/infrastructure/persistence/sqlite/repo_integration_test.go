package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"shopping/internal/domain/products"
	"shopping/internal/migrator"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	return db
}

func TestRepo_ListProducts_FilteringAndPaging_NoNamedArgError(t *testing.T) {
	db := openTestDB(t)
	if err := migrator.Up(db); err != nil {
		t.Fatalf("migrator.Up: %v", err)
	}

	r := NewRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Regression: empty group filter used to trigger sqlite/sqlc placeholder issues
	// ("missing named argument \"6\"") due to generated LIMIT ?6 / OFFSET ?5.
	_, err := r.ListProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: false,
		NameQuery:        "",
		GroupIDs:         nil,
		Limit:            products.MaxProductsPageSize,
		Offset:           0,
	})
	if err != nil {
		t.Fatalf("ListProducts (empty groups): %v", err)
	}
	_, err = r.CountProducts(ctx, products.ProductFilter{
		OnlyMissingOrLow: false,
		NameQuery:        "",
		GroupIDs:         nil,
	})
	if err != nil {
		t.Fatalf("CountProducts (empty groups): %v", err)
	}

	// Verify group + name filters work (multi-select semantics).
	groups, err := r.ListGroups(ctx)
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	var warzywaID, makiID products.GroupID
	for _, g := range groups {
		switch g.Name {
		case "warzywa":
			warzywaID = g.ID
		case "mąki":
			makiID = g.ID
		}
	}
	if warzywaID == 0 || makiID == 0 {
		t.Fatalf("expected seeded groups warzywa and mąki, got warzywa=%d mąki=%d", warzywaID, makiID)
	}

	all, err := r.ListProducts(ctx, products.ProductFilter{
		Limit: products.MaxProductsPageSize,
	})
	if err != nil {
		t.Fatalf("ListProducts (unfiltered): %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("expected seeded products to exist")
	}
	seedFound := false
	for _, p := range all {
		if p.Name == "mąka razowa" {
			seedFound = true
			break
		}
	}
	if !seedFound {
		t.Fatalf("expected to find seeded product %q in unfiltered list", "mąka razowa")
	}

	byGroups, err := r.ListProducts(ctx, products.ProductFilter{
		GroupIDs: []products.GroupID{warzywaID, makiID},
		Limit:    products.MaxProductsPageSize,
	})
	if err != nil {
		t.Fatalf("ListProducts (groups): %v", err)
	}
	if len(byGroups) == 0 {
		t.Fatalf("expected some results for groups filter")
	}

	byName, err := r.ListProducts(ctx, products.ProductFilter{
		NameQuery: "razowa",
		Limit:     products.MaxProductsPageSize,
	})
	if err != nil {
		t.Fatalf("ListProducts (name): %v", err)
	}
	if len(byName) == 0 {
		t.Fatalf("expected some results for name filter")
	}

	list, err := r.ListProducts(ctx, products.ProductFilter{
		NameQuery: "razowa",
		GroupIDs:  []products.GroupID{warzywaID, makiID},
		Limit:     products.MaxProductsPageSize,
	})
	if err != nil {
		t.Fatalf("ListProducts (name+groups): %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected some results for name+groups filter")
	}
	foundFlour := false
	for _, p := range list {
		if p.Name == "mąka razowa" {
			foundFlour = true
			break
		}
	}
	if !foundFlour {
		t.Fatalf("expected to find seeded product %q in filtered results", "mąka razowa")
	}

	// Verify paging clamp behavior doesn't error.
	_, err = r.ListProducts(ctx, products.ProductFilter{
		Limit:  9999,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("ListProducts (limit clamp): %v", err)
	}
}
