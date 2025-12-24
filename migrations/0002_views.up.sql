CREATE VIEW IF NOT EXISTS v_groups AS
SELECT id, name
FROM groups;

CREATE VIEW IF NOT EXISTS v_products AS
SELECT
  p.id,
  p.name,
  p.group_id,
  g.name AS group_name,
  p.quantity,
  p.min_quantity,
  p.missing,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;

