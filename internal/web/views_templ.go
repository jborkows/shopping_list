package web

import (
	"context"
	"io"
	"net/url"
	"shopping/internal/domain/products"
	"strconv"
	"strings"
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
			`<link rel="stylesheet" href="/static/app.css?v=2" />`,
			`<script src="`, templ.EscapeString(data.Base.HTMXSrc), `"></script>`,
			`<script>`,
			`document.addEventListener('input', function (e) {`,
			`  var t = e.target;`,
			`  if (!t || t.getAttribute('data-ms-filter') !== 'groups') return;`,
			`  var root = t.closest('.multiselect');`,
			`  if (!root) return;`,
			`  var q = (t.value || '').toLowerCase().trim();`,
			`  var opts = root.querySelectorAll('.ms-option');`,
			`  for (var i = 0; i < opts.length; i++) {`,
			`    var text = (opts[i].textContent || '').toLowerCase();`,
			`    opts[i].style.display = (q === '' || text.indexOf(q) !== -1) ? '' : 'none';`,
			`  }`,
			`});`,
			`function updateGroupSummary(root) {`,
			`  if (!root) return;`,
			`  var countEl = root.querySelector('.ms-count');`,
			`  if (!countEl) return;`,
			`  var checked = root.querySelectorAll('input[type=checkbox][name=group_id]:checked');`,
			`  countEl.textContent = String(checked.length);`,
			`}`,
			`document.addEventListener('change', function (e) {`,
			`  var t = e.target;`,
			`  if (!t || t.name !== 'group_id' || t.type !== 'checkbox') return;`,
			`  var root = t.closest('.multiselect');`,
			`  updateGroupSummary(root);`,
			`});`,
			`document.addEventListener('DOMContentLoaded', function () {`,
			`  var roots = document.querySelectorAll('details.multiselect');`,
			`  for (var i = 0; i < roots.length; i++) updateGroupSummary(roots[i]);`,
			`});`,
			`document.addEventListener('click', function (e) {`,
			`  var target = e.target;`,
			`  var open = document.querySelectorAll('details.multiselect[open]');`,
			`  for (var i = 0; i < open.length; i++) {`,
			`    if (open[i].contains(target)) continue;`,
			`    open[i].removeAttribute('open');`,
			`  }`,
			`});`,
			`</script>`,
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

func AdminPage(data adminPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := write(w,
			"<!doctype html>",
			`<html lang="en">`,
			`<head>`,
			`<meta charset="utf-8" />`,
			`<meta name="viewport" content="width=device-width, initial-scale=1" />`,
			`<title>`, templ.EscapeString(data.Base.Title), `</title>`,
			`<link rel="stylesheet" href="/static/app.css?v=2" />`,
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

		if err := renderAdminContent(w); err != nil {
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

func ShoppingListPage(data shoppingListPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := write(w,
			"<!doctype html>",
			`<html lang="en">`,
			`<head>`,
			`<meta charset="utf-8" />`,
			`<meta name="viewport" content="width=device-width, initial-scale=1" />`,
			`<title>`, templ.EscapeString(data.Base.Title), `</title>`,
			`<link rel="stylesheet" href="/static/app.css?v=2" />`,
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

		if err := write(w, `<div class="row"><h1>Lista zakupów</h1><div class="spacer"></div></div>`); err != nil {
			return err
		}
		if err := renderShoppingListCard(w, shoppingListData{Items: data.Items}); err != nil {
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

func ShoppingListCard(data shoppingListData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return renderShoppingListCard(w, data)
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
	if err := write(w, `<nav class="nav"><a href="/products">Zapasy</a><a href="/shopping-list">Lista zakupów</a>`); err != nil {
		return err
	}
	if base.IsAdmin {
		if err := write(w, `<a href="/admin">Admin</a>`); err != nil {
			return err
		}
	}
	if err := write(w, `</nav>`); err != nil {
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
		`<h1>Zapasy</h1>`,
		`<div class="spacer"></div>`,
		`<a class="`, allClass, `" href="/products`, productsListQS(false, data.NameQuery, data.SelectedGroupIDs, 1), `">Wszystko</a>`,
		`<a class="`, missingClass, `" href="/products`, productsListQS(true, data.NameQuery, data.SelectedGroupIDs, 1), `">Braki / niski stan</a>`,
		`</div>`,
	); err != nil {
		return err
	}

	listQS := productsListQS(data.OnlyMissing, data.NameQuery, data.SelectedGroupIDs, data.Page)

	if err := write(w, `<section class="card"><h2>Dodaj produkt</h2>`); err != nil {
		return err
	}
	if err := write(w,
		`<form hx-post="/products`, listQS, `" hx-target="#products-list" hx-swap="outerHTML">`,
		`<div class="form-row">`,
		`<input name="name" placeholder="np. marchewka" required />`,
		`<input class="num-input" type="number" step="0.1" min="0" name="quantity" value="0" />`,
		`<select name="unit">`, unitOptions(string(products.UnitPiece)), `</select>`,
		`<input class="num-input" type="number" step="0.1" min="0" name="min_quantity" value="0" />`,
		`<select name="group_id">`,
		`<option value="">(brak grupy)</option>`,
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

	if err := write(w, `<h2>Utwórz grupę</h2>`); err != nil {
		return err
	}
	if err := write(w,
		`<form hx-post="/groups`, listQS, `" hx-swap="none">`,
		`<div class="form-row">`,
		`<input name="name" placeholder="np. warzywa" required />`,
		`<button type="submit">Utwórz</button>`,
		`</div></form>`,
		`</section>`,
	); err != nil {
		return err
	}

	return renderProductsList(w, productsListData{
		Groups:           data.Groups,
		Products:         data.Products,
		OnlyMissing:      data.OnlyMissing,
		NameQuery:        data.NameQuery,
		SelectedGroupIDs: data.SelectedGroupIDs,
		Page:             data.Page,
		TotalPages:       data.TotalPages,
		Total:            data.Total,
	})
}

func renderAdminContent(w io.Writer) error {
	if err := write(w, `<div class="row"><h1>Administracja</h1><div class="spacer"></div></div>`); err != nil {
		return err
	}
	return write(w,
		`<section class="card">`,
		`<h2>Utrzymanie</h2>`,
		`<div class="form-row">`,
		`<button hx-post="/admin/db/optimize" hx-target="#admin-status" hx-swap="innerHTML">Optymalizuj bazę</button>`,
		`<span class="muted" id="admin-status"></span>`,
		`</div>`,
		`</section>`,
	)
}

func renderShoppingListCard(w io.Writer, data shoppingListData) error {
	if err := write(w,
		`<section class="card" id="shopping-list">`,
		`<h2>Lista zakupów</h2>`,
		`<form hx-post="/shopping-list" hx-target="#shopping-list" hx-swap="outerHTML">`,
		`<div class="form-row">`,
		`<input name="name" placeholder="np. mleko" required autocomplete="off" hx-get="/partials/product-suggestions" hx-target="#product-suggestions" hx-swap="innerHTML" hx-trigger="keyup changed delay:200ms" />`,
		`<input class="num-input" type="number" step="0.1" min="0.1" name="quantity" value="1" required />`,
		`<select name="unit">`, unitOptions(string(products.UnitPiece)), `</select>`,
		`<button type="submit">Dodaj</button>`,
		`</div>`,
		`<div id="product-suggestions"></div>`,
		`</form>`,
	); err != nil {
		return err
	}

	if len(data.Items) == 0 {
		if err := write(w, `<p class="muted">Brak produktów na liście. Użyj „Na listę” w zapasach albo dodaj ręcznie powyżej.</p>`); err != nil {
			return err
		}
		return write(w, `</section>`)
	}

	if err := write(w, `<div class="sl-items">`); err != nil {
		return err
	}
	for _, item := range data.Items {
		id := strconv.FormatInt(int64(item.ID), 10)
		nextDone := "1"
		nameClass := ""
		doneBtnClass := "icon-btn sl-done sl-done-off"
		doneInner := `<span class="sl-done-empty" aria-hidden="true"></span>`
		if item.Done {
			nextDone = "0"
			nameClass = ` class="done"`
			doneBtnClass = "icon-btn sl-done sl-done-on"
			doneInner = "✅"
		}

		iconKey := strings.TrimSpace(item.IconKey)
		iconSrc := ""
		if iconKey != "" {
			iconSrc = "/static/icons/" + iconKey + ".svg"
		} else {
			iconSrc = "/icons/auto?name=" + url.QueryEscape(item.Name)
		}
		group := strings.TrimSpace(item.GroupName)

		if err := write(w, `<div class="sl-item">`); err != nil {
			return err
		}
		if err := write(w, `<div class="sl-top">`); err != nil {
			return err
		}
		if err := write(w, `<div class="sl-name`, nameClass, `">`, templ.EscapeString(item.Name)); err != nil {
			return err
		}
		if group != "" {
			if err := write(w, ` <span class="muted">(`, templ.EscapeString(group), `)</span>`); err != nil {
				return err
			}
		}
		if err := write(w, `</div></div>`); err != nil {
			return err
		}

		if err := write(w, `<div class="sl-mid">`); err != nil {
			return err
		}
		if err := write(w, `<img class="sl-icon" width="32" height="32" src="`, templ.EscapeString(iconSrc), `" alt="" />`); err != nil {
			return err
		}
		if err := write(w,
			`<form class="sl-qty" hx-patch="/shopping-list/items/`, id, `" hx-target="#shopping-list" hx-swap="outerHTML" hx-trigger="change">`,
			`<input class="num-input" type="number" step="0.1" min="0.1" name="quantity" value="`, formatQty(item.Quantity), `" />`,
			`<select name="unit">`, unitOptions(string(item.Unit)), `</select>`,
			`</form>`,
		); err != nil {
			return err
		}
		if err := write(w,
			`<form class="sl-done-form" hx-patch="/shopping-list/items/`, id, `" hx-target="#shopping-list" hx-swap="outerHTML">`,
			`<input type="hidden" name="done" value="`, nextDone, `" />`,
			`<button class="`, doneBtnClass, `" type="submit" aria-label="Zrobione">`, doneInner, `</button>`,
			`</form>`,
		); err != nil {
			return err
		}

		if err := write(w, `<div class="sl-right">`); err != nil {
			return err
		}
		if item.ProductID == nil {
			if err := write(w,
				`<form hx-post="/shopping-list/items/`, id, `/product" hx-swap="none">`,
				`<button class="icon-btn" type="submit" aria-label="Dodaj do zapasów"><img class="icon" src="/static/icons/plus.svg" alt="" /></button>`,
				`</form>`,
			); err != nil {
				return err
			}
		}
		if err := write(w,
			`<form class="sl-delete" hx-delete="/shopping-list/items/`, id, `" hx-target="#shopping-list" hx-swap="outerHTML">`,
			`<button class="icon-btn secondary" type="submit" aria-label="Usuń">🗑️</button>`,
			`</form>`,
			`</div>`, // sl-right
			`</div>`, // sl-mid
			`</div>`, // sl-item
		); err != nil {
			return err
		}
	}
	if err := write(w, `</div>`); err != nil {
		return err
	}
	return write(w, `</section>`)
}

func renderProductsList(w io.Writer, data productsListData) error {
	title := "Wszystkie produkty"
	if data.OnlyMissing {
		title = "Braki / niski stan"
	}

	listQS := productsListQS(data.OnlyMissing, data.NameQuery, data.SelectedGroupIDs, data.Page)

	if err := write(w,
		`<section class="card" id="products-list">`,
		`<div class="row">`,
		`<h2>`, templ.EscapeString(title), `</h2>`,
		`<div class="spacer"></div>`,
		`<span class="muted">Wyniki: `, strconv.FormatInt(data.Total, 10), `</span>`,
		`</div>`,
		`<form id="products-filters" class="filters" method="get" action="/products" hx-get="/partials/products" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="submit, keyup changed delay:250ms from:#products-q">`,
		`<input type="hidden" name="missing" value="`, boolToQuery(data.OnlyMissing), `" />`,
		`<input type="hidden" name="page" value="1" />`,
		`<div class="form-row">`,
		`<input id="products-q" name="q" value="`, templ.EscapeString(data.NameQuery), `" placeholder="Szukaj po nazwie…" />`,
		`<details class="multiselect">`,
		`<summary>Grupy (<span class="ms-count">`, strconv.FormatInt(int64(len(data.SelectedGroupIDs)), 10), `</span>)</summary>`,
		`<div class="multiselect-popover" role="listbox" aria-label="Grupy">`,
		`<input class="ms-search" type="search" placeholder="Szukaj grup…" autocomplete="off" data-ms-filter="groups" />`,
	); err != nil {
		return err
	}

	selected := make(map[products.GroupID]struct{}, len(data.SelectedGroupIDs))
	for _, gid := range data.SelectedGroupIDs {
		selected[gid] = struct{}{}
	}
	for _, g := range data.Groups {
		_, isSelected := selected[g.ID]
		if err := write(w, `<label class="ms-option"><input type="checkbox" name="group_id" value="`, strconv.FormatInt(int64(g.ID), 10), `"`); err != nil {
			return err
		}
		if isSelected {
			if err := write(w, ` checked`); err != nil {
				return err
			}
		}
		if err := write(w, ` />`, templ.EscapeString(g.Name), `</label>`); err != nil {
			return err
		}
	}
	if err := write(w,
		`</div>`,
		`</details>`,
		`<button class="secondary" type="submit">Filtruj</button>`,
		`</div>`,
		`</form>`,
	); err != nil {
		return err
	}

	if data.TotalPages > 1 {
		prevPage := data.Page - 1
		nextPage := data.Page + 1
		if prevPage < 1 {
			prevPage = 1
		}
		if nextPage > data.TotalPages {
			nextPage = data.TotalPages
		}

		if err := write(w, `<div class="row pager">`); err != nil {
			return err
		}
		if data.Page > 1 {
			prevQS := productsListQS(data.OnlyMissing, data.NameQuery, data.SelectedGroupIDs, prevPage)
			if err := write(w, `<a class="button secondary" href="/products`, prevQS, `" hx-get="/partials/products`, prevQS, `" hx-target="#products-list" hx-swap="outerHTML">Poprzednia</a>`); err != nil {
				return err
			}
		} else {
			if err := write(w, `<span class="button secondary disabled">Poprzednia</span>`); err != nil {
				return err
			}
		}
		if err := write(w, `<span class="muted">Strona `, strconv.FormatInt(data.Page, 10), ` z `, strconv.FormatInt(data.TotalPages, 10), `</span>`); err != nil {
			return err
		}
		if data.Page < data.TotalPages {
			nextQS := productsListQS(data.OnlyMissing, data.NameQuery, data.SelectedGroupIDs, nextPage)
			if err := write(w, `<a class="button secondary" href="/products`, nextQS, `" hx-get="/partials/products`, nextQS, `" hx-target="#products-list" hx-swap="outerHTML">Następna</a>`); err != nil {
				return err
			}
		} else {
			if err := write(w, `<span class="button secondary disabled">Następna</span>`); err != nil {
				return err
			}
		}
		if err := write(w, `</div>`); err != nil {
			return err
		}
	}

	if err := write(w,
		`<table class="table">`,
		`<thead><tr>`,
		`<th>Nazwa</th><th>Grupa</th><th class="num">Ilość</th><th class="num">Min</th><th>Brak</th><th>Aktualizacja</th><th></th>`,
		`</tr></thead>`,
		`<tbody>`,
	); err != nil {
		return err
	}

	for _, p := range data.Products {
		rowClass := ""
		if p.Missing || p.Quantity <= p.MinQuantity {
			rowClass = ` class="warn"`
		}

		if err := write(w, `<tr`, rowClass, `>`); err != nil {
			return err
		}
		iconKey := strings.TrimSpace(p.IconKey)
		if iconKey == "" {
			iconKey = "cart"
		}
		if err := write(w, `<td><div class="product-cell"><img class="product-icon" width="48" height="48" src="/static/icons/`, templ.EscapeString(iconKey), `.svg" alt="" /><span class="product-name">`, templ.EscapeString(p.Name), `</span></div></td>`); err != nil {
			return err
		}

		// Group select.
		if err := write(w, `<td>`); err != nil {
			return err
		}
		if err := write(w, `<form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/group`, listQS, `" hx-target=\"#products-list\" hx-swap=\"outerHTML\" hx-trigger=\"change\">`); err != nil {
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
		if err := write(w, `>(brak grupy)</option>`); err != nil {
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

		// Quantity + unit.
		if err := write(w, `<td class="num">`); err != nil {
			return err
		}
		if err := write(w, `<div class="form-row">`); err != nil {
			return err
		}
		if err := write(w, `<form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/qty`, listQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<input class="num-input" type="number" step="0.1" min="0" name="quantity" value="`, formatQty(p.Quantity), `" />`); err != nil {
			return err
		}
		if err := write(w, `</form>`); err != nil {
			return err
		}
		if err := write(w, `<form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/unit`, listQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<select name="unit">`, unitOptions(string(p.Unit)), `</select></form>`); err != nil {
			return err
		}
		if err := write(w, `</div></td>`); err != nil {
			return err
		}

		// Min.
		if err := write(w, `<td class="num"><div class="form-row">`); err != nil {
			return err
		}
		if err := write(w, `<form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/min`, listQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
			return err
		}
		if err := write(w, `<input class="num-input" type="number" step="0.1" min="0" name="min_quantity" value="`, formatQty(p.MinQuantity), `" />`); err != nil {
			return err
		}
		if err := write(w, `</form><span class="muted">`, templ.EscapeString(string(p.Unit)), `</span></div></td>`); err != nil {
			return err
		}

		// Missing checkbox.
		if err := write(w, `<td><form hx-post="/products/`, strconv.FormatInt(int64(p.ID), 10), `/missing`, listQS, `" hx-target="#products-list" hx-swap="outerHTML" hx-trigger="change">`); err != nil {
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

		if err := write(w, `<td>`); err != nil {
			return err
		}
		if err := write(w,
			`<form hx-post="/shopping-list/from-product/`, strconv.FormatInt(int64(p.ID), 10), `" hx-swap="none">`,
			`<button type="submit">Na listę</button>`,
			`</form>`,
			`</td>`,
		); err != nil {
			return err
		}

		if err := write(w, `</tr>`); err != nil {
			return err
		}
	}

	if len(data.Products) == 0 {
		if err := write(w, `<tr><td colspan="7" class="muted">Brak produktów.</td></tr>`); err != nil {
			return err
		}
	}

	return write(w, `</tbody></table></section>`)
}

func boolToQuery(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func productsListQS(onlyMissing bool, nameQuery string, groupIDs []products.GroupID, page int64) string {
	values := url.Values{}
	if onlyMissing {
		values.Set("missing", "1")
	}
	nameQuery = strings.TrimSpace(nameQuery)
	if nameQuery != "" {
		values.Set("q", nameQuery)
	}
	for _, gid := range groupIDs {
		values.Add("group_id", strconv.FormatInt(int64(gid), 10))
	}
	if page > 1 {
		values.Set("page", strconv.FormatInt(page, 10))
	}
	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}

func write(w io.Writer, parts ...string) error {
	for _, p := range parts {
		if _, err := io.WriteString(w, p); err != nil {
			return err
		}
	}
	return nil
}

func formatQty(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func unitOptions(selected string) string {
	selected = strings.TrimSpace(selected)
	if selected == "" {
		selected = string(products.UnitPiece)
	}
	opts := []struct {
		Value string
		Label string
	}{
		{Value: string(products.UnitKG), Label: string(products.UnitKG)},
		{Value: string(products.UnitLiter), Label: string(products.UnitLiter)},
		{Value: string(products.UnitPiece), Label: string(products.UnitPiece)},
		{Value: string(products.UnitGram), Label: string(products.UnitGram)},
	}
	var b strings.Builder
	for _, o := range opts {
		if o.Value == selected {
			b.WriteString(`<option value="` + o.Value + `" selected>` + o.Label + `</option>`)
			continue
		}
		b.WriteString(`<option value="` + o.Value + `">` + o.Label + `</option>`)
	}
	return b.String()
}
