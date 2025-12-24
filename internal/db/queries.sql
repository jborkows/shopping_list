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
  p.group_id,
  p.group_name,
  p.quantity,
  p.min_quantity,
  p.missing,
  p.updated_at
FROM v_products p
ORDER BY p.name;

-- name: ListProductsMissingOrLow :many
SELECT
  p.id,
  p.name,
  p.group_id,
  p.group_name,
  p.quantity,
  p.min_quantity,
  p.missing,
  p.updated_at
FROM v_products p
WHERE p.missing = 1 OR p.quantity <= p.min_quantity
ORDER BY p.name;

-- name: CreateProduct :one
INSERT INTO products(name, group_id, quantity, min_quantity, missing, created_at, updated_at)
VALUES (?, sqlc.narg('group_id'), 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id;

-- name: SetProductQuantity :exec
UPDATE products
SET quantity = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductMinQuantity :exec
UPDATE products
SET min_quantity = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductMissing :exec
UPDATE products
SET missing = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetProductGroup :exec
UPDATE products
SET group_id = sqlc.narg('group_id'), updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: PragmaOptimize :exec
PRAGMA optimize;
