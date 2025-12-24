package web

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func ProductsPage(data productsPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := write(w,
			"<!doctype html>",
			`<html lang="en">`,
			`<head>`,
			`<meta charset="utf-8" />`,
			`<meta name="viewport" content="width=device-width, initial-scale=1" />`,
			`<title>`, templ.EscapeString(data.Base.Title), `</title>`,
			`<link rel="stylesheet" href="/static/app.css" />`,
			`<script src="`, templ.EscapeString(data.Base.HTMXSrc), `"></script>`,
			`</head>`,
			`<body>`,
		); err != nil {
			return err
		}

		if err := renderTopbar(w, data.Base); err != nil {
			return err
		}

		if err := write(w, `<main class="container">`); err != nil {
			return err
		}

		if err := renderProductsContent(w, data); err != nil {
			return err
		}

		if err := write(w, `</main>`); err != nil {
			return err
		}

		if err := write(w, `<footer class="container footer"><span>© `, strconv.Itoa(time.Now().Year()), `</span></footer>`); err != nil {
			return err
		}

		return write(w, `</body></html>`)
	})
}

func ProductsList(data productsListData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return renderProductsList(w, data)
	})
}

func renderTopbar(w io.Writer, base baseData) error {
	if err := write(w, `<header class="topbar"><div class="container topbar-row">`); err != nil {
		return err
	}
	if err := write(w, `<div class="brand"><a href="/products">shopping</a></div>`); err != nil {
		return err
	}
	if err := write(w, `<nav class="nav"><a href="/products">Products</a></nav>`); err != nil {
		return err
	}
	if err := write(w, `<div class="spacer"></div>`); err != nil {
		return err
	}

	if base.User != nil {
		name := base.User.Email
		if base.User.Name != "" {
			name = base.User.Name
		}
		if err := write(w, `<div class="user"><span class="user-name">`, templ.EscapeString(name), `</span>`); err != nil {
			return err
		}
		if err := write(w, `<form method="post" action="/logout"><button class="link" type="submit">Logout</button></form></div>`); err != nil {
			return err
		}
	}

	return write(w, `</div></header>`)
}

func renderProductsContent(w io.Writer, data productsPageData) error {
	allClass := `button`
	missingClass := `button`
	if data.OnlyMissing {
		allClass = `button secondary`
	} else {
		missingClass = `button secondary`
	}

	if err := write(w,
		`<div class="row">`,
		`<h1>Supplies</h1>`,
		`<div class="spacer"></div>`,
		`<a class="`, allClass, `" href="/products">All</a>`,
		`<a class="`, missingClass, `" href="/products?missing=1">Missing / Low</a>`,
		`</div>`,
	); err != nil {
		return err
	}

	missingQS := ""
	if data.OnlyMissing {
		missingQS = "?missing=1"
	}

	if err := write(w, `<section class="card"><h2>Add product</h2>`); err != nil {
		return err
	}
	if err := write(w,
		`<form hx-post="/products`, missingQS, `" hx-target="#products-list" hx-swap="outerHTML">`,
		`<div class="form-row">`,
		`<input name="name" placeholder="e.g. carrots" required />`,
		`<select name="group_id">`,
		`<option value="">(no group)</option>`,
	); err != nil {
		return err
	}
	for _, g := range data.Groups {
		if err := write(w, `<option value="`, strconv.FormatInt(int64(g.ID), 10), `">`, templ.EscapeString(g.Name), `</option>`); err != nil {
			return err
		}
	}
	if err := write(w,
		`</select>`,
		`<button type="submit">Add</button>`,
		`</div></form>`,
	); err != nil {
		return err
	}

	if err := write(w, `<h2>Create group</h2>`); err != nil {
		return err
	}
	if err := write(w,
		`<form hx-post="/groups`, missingQS, `" hx-swap="none">`,
		`<div class="form-row">`,
		`<input name="name" placeholder="e.g. vegetables" required />`,
		`<button type="submit">Create</button>`,
		`</div></form>`,
		`</section>`,
	); err != nil {
		return err
	}

	if data.Base.IsAdmin {
		if err := write(w,
			`<section class="card">`,
			`<h2>Admin</h2>`,
			`<div class="form-row">`,
			`<button hx-post="/admin/db/optimize" hx-target="#admin-status" hx-swap="innerHTML">Optimize database</button>`,
			`<span class="muted" id="admin-status"></span>`,
			`</div>`,
			`</section>`,
		); err != nil {
			return err
		}
	}

	return renderProductsList(w, productsListData{
		Groups:      data.Groups,
		Products:    data.Products,
		OnlyMissing: data.OnlyMissing,
	})
}

func renderProductsList(w io.Writer, data productsListData) error {
	title := "All products"
	if data.OnlyMissing {
		title = "Missing / Low"
	}
	if err := write(w,
		`<section class="card" id="products-list">`,
		`<h2>`, templ.EscapeString(title), `</h2>`,
		`<table class="table">`,
		`<thead><tr>`,
		`<th>Name</th><th>Group</th><th class="num">Qty</th><th class="num">Min</th><th>Missing</th><th>Updated</th>`,
		`</tr></thead>`,
		`<tbody>`,
	); err != nil {
		return err
	}

	missingQS := ""
	if data.OnlyMissing {
		missingQS = "?missing=1"
	}

	for _, p := range data.Products {
		rowClass := ""
		if p.Missing || p.Quantity <= p.MinQuantity {
			rowClass = ` class="warn"`
		}

		if err := write(w, `<tr`, rowClass, `>`); err != nil {
			return err
		}
		if err := write(w, `<td>`, templ.EscapeString(p.Name), `</td>`); err != nil {
			return err
		}

		// Group select.
		if err := write(w, `<td>`); err != nil {
			return err
		}
		if err := write(w, `<form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/group`, missingQS, `" hx-target=\"#products-list\" hx-swap=\"outerHTML\" hx-trigger=\"change\">`); err != nil {
			return err
		}
		if err := write(w, `<select name="group_id">`); err != nil {
			return err
		}
		if err := write(w, `<option value=""`); err != nil {
			return err
		}
		if p.GroupID == nil {
			if err := write(w, ` selected`); err != nil {
				return err
			}
		}
		if err := write(w, `>(no group)</option>`); err != nil {
			return err
		}
		for _, g := range data.Groups {
			selected := p.GroupID != nil && *p.GroupID == g.ID
			if err := write(w, `<option value="`, strconv.FormatInt(int64(g.ID), 10), `"`); err != nil {
				return err
			}
			if selected {
				if err := write(w, ` selected`); err != nil {
					return err
				}
			}
			if err := write(w, `>`, templ.EscapeString(g.Name), `</option>`); err != nil {
				return err
			}
		}
		if err := write(w, `</select></form></td>`); err != nil {
			return err
		}

		// Quantity.
		if err := write(w, `<td class="num"><form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/qty`, missingQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<input class="num-input" type="number" name="quantity" value="`, strconv.Itoa(p.Quantity), `" />`); err != nil {
			return err
		}
		if err := write(w, `</form></td>`); err != nil {
			return err
		}

		// Min.
		if err := write(w, `<td class="num"><form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/min`, missingQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<input class="num-input" type="number" name="min_quantity" value="`, strconv.Itoa(p.MinQuantity), `" />`); err != nil {
			return err
		}
		if err := write(w, `</form></td>`); err != nil {
			return err
		}

		// Missing checkbox.
		if err := write(w, `<td><form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/missing`, missingQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<input type="checkbox" name="missing"`); err != nil {
			return err
		}
		if p.Missing {
			if err := write(w, ` checked`); err != nil {
				return err
			}
		}
		if err := write(w, ` /></form></td>`); err != nil {
			return err
		}

		if err := write(w, `<td class="muted">`, templ.EscapeString(p.UpdatedAt.Format("2006-01-02 15:04")), `</td>`); err != nil {
			return err
		}

		if err := write(w, `</tr>`); err != nil {
			return err
		}
	}

	if len(data.Products) == 0 {
		if err := write(w, `<tr><td colspan="6" class="muted">No products yet.</td></tr>`); err != nil {
			return err
		}
	}

	return write(w, `</tbody></table></section>`)
}

func write(w io.Writer, parts ...string) error {
	for _, p := range parts {
		if _, err := io.WriteString(w, p); err != nil {
			return err
		}
	}
	return nil
}
