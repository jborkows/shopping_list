-- name: ListGroups :many
SELECT id, name
FROM v_groups
ORDER BY name;

-- name: CreateGroup :one
INSERT INTO groups(name) VALUES (?)
RETURNING id;

-- name: ListProductsAll :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.missing,
  p.updated_at
FROM v_products p
ORDER BY p.name;

-- name: ListProductsMissingOrLow :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.missing,
  p.updated_at
FROM v_products p
WHERE p.missing = 1 OR p.quantity_value <= p.min_quantity_value
ORDER BY p.name;

-- name: ListProductsFiltered :many
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  p.group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.missing,
  p.updated_at
FROM v_products p
WHERE
  (? = 0 OR p.missing = 1 OR p.quantity_value <= p.min_quantity_value)
  AND (? = '' OR lower(p.name) LIKE '%' || lower(?) || '%')
  AND (? = 0 OR p.group_id IN (sqlc.slice('group_ids')))
ORDER BY lower(p.name)
LIMIT ?
OFFSET ?;

-- name: CountProductsFiltered :one
SELECT COUNT(*)
FROM v_products p
WHERE
  (? = 0 OR p.missing = 1 OR p.quantity_value <= p.min_quantity_value)
  AND (? = '' OR lower(p.name) LIKE '%' || lower(?) || '%')
  AND (? = 0 OR p.group_id IN (sqlc.slice('group_ids')));

-- name: CreateProduct :one
INSERT INTO products(name, icon_key, group_id, quantity_value, quantity_unit, min_quantity_value, missing, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id;

-- name: SetProductQuantity :exec
UPDATE products
SET quantity_value = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductMinQuantity :exec
UPDATE products
SET min_quantity_value = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductUnit :exec
UPDATE products
SET quantity_unit = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductMissing :exec
UPDATE products
SET missing = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductGroup :exec
UPDATE products
SET group_id = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListShoppingListItems :many
SELECT
  sli.id,
  sli.product_id,
  sli.name,
  COALESCE(p.icon_key, '') AS icon_key,
  COALESCE(g.name, '') AS group_name,
  sli.quantity_value,
  sli.quantity_unit,
  sli.done,
  sli.created_at
FROM shopping_list_items sli
LEFT JOIN products p ON p.id = sli.product_id
LEFT JOIN groups g ON g.id = p.group_id
ORDER BY sli.created_at DESC, sli.id DESC;

-- name: GetShoppingListItem :one
SELECT
  sli.id,
  sli.product_id,
  sli.name,
  COALESCE(p.icon_key, '') AS icon_key,
  COALESCE(g.name, '') AS group_name,
  sli.quantity_value,
  sli.quantity_unit,
  sli.done,
  sli.created_at
FROM shopping_list_items sli
LEFT JOIN products p ON p.id = sli.product_id
LEFT JOIN groups g ON g.id = p.group_id
WHERE sli.id = ?;

-- name: AddShoppingListItemByName :exec
INSERT OR IGNORE INTO shopping_list_items(product_id, name, quantity_value, quantity_unit, done)
VALUES (NULL, ?, ?, ?, 0);

-- name: AddShoppingListItemByProductID :exec
INSERT OR IGNORE INTO shopping_list_items(product_id, name, quantity_value, quantity_unit, done)
SELECT p.id, p.name, 1, p.quantity_unit, 0
FROM products p
WHERE p.id = ?;

-- name: SetShoppingListItemDone :exec
UPDATE shopping_list_items
SET done = ?
WHERE id = ?;

-- name: SetShoppingListItemQuantity :exec
UPDATE shopping_list_items
SET quantity_value = ?, quantity_unit = ?
WHERE id = ?;

-- name: DeleteShoppingListItem :exec
DELETE FROM shopping_list_items
WHERE id = ?;

-- name: LinkShoppingListItemToProduct :exec
UPDATE shopping_list_items
SET product_id = ?, name = ?
WHERE id = ?;

-- name: FindProductIDByName :one
SELECT id
FROM products
WHERE lower(name) = lower(?)
LIMIT 1;

-- name: SuggestProductsByName :many
SELECT
  id,
  name,
  icon_key,
  quantity_unit
FROM products
WHERE lower(name) LIKE '%' || lower(?) || '%'
ORDER BY lower(name)
LIMIT ?;

-- name: ResolveProductIconKeyByName :one
SELECT icon_key
FROM product_icon_rules
WHERE lower(?) LIKE '%' || lower(match_substring) || '%'
ORDER BY priority DESC, id DESC
LIMIT 1;
