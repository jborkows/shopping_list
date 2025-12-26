CREATE TABLE IF NOT EXISTS groups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  group_id INTEGER NULL REFERENCES groups(id) ON DELETE SET NULL,
  icon_key TEXT NOT NULL DEFAULT 'cart',
  quantity_value REAL NOT NULL DEFAULT 0,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  min_quantity_value REAL NOT NULL DEFAULT 0,
  missing INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_group_id ON products(group_id);

CREATE TABLE IF NOT EXISTS shopping_list_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  quantity_value REAL NOT NULL DEFAULT 1,
  quantity_unit TEXT NOT NULL DEFAULT 'sztuk',
  done INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_shopping_list_open_product
  ON shopping_list_items(product_id)
  WHERE done = 0 AND product_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_shopping_list_open_name
  ON shopping_list_items(lower(name))
  WHERE done = 0 AND product_id IS NULL;

CREATE TABLE IF NOT EXISTS product_icon_rules (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  match_substring TEXT NOT NULL,
  icon_key TEXT NOT NULL,
  priority INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_product_icon_rules_priority
  ON product_icon_rules(priority DESC, id DESC);

CREATE VIEW IF NOT EXISTS v_groups AS
SELECT id, name
FROM groups;

CREATE VIEW IF NOT EXISTS v_products AS
SELECT
  p.id,
  p.name,
  p.icon_key,
  p.group_id,
  g.name AS group_name,
  p.quantity_value,
  p.quantity_unit,
  p.min_quantity_value,
  p.missing,
  p.updated_at
FROM products p
LEFT JOIN groups g ON g.id = p.group_id;

-- Seed groups.
INSERT OR IGNORE INTO groups(name) VALUES
  ('warzywa'),
  ('owoce'),
  ('jajka'),
  ('nabiał'),
  ('spożywcze'),
  ('mąki'),
  ('domowe'),
  ('chemia'),
  ('makarony'),
  ('ryż'),
  ('mięso');

-- Seed icon rules (used on new product creation).
INSERT INTO product_icon_rules(match_substring, icon_key, priority)
VALUES
  ('marchew', 'carrot', 100),
  ('ziemni', 'potato', 100),
  ('mąka', 'flour', 100),
  ('maka', 'flour', 90),
  ('jaj', 'eggs', 100),
  ('mleko', 'milk', 100),
  ('wołow', 'cow', 130),
  ('wołowe', 'cow', 130),
  ('kurcz', 'chicken', 130),
  ('wieprz', 'pig', 130),
  ('schab', 'pig', 135),
  ('karków', 'pig', 135),
  ('boczek', 'bacon', 110),
  ('papier toalet', 'toilet-paper', 100),
  ('indyk', 'turkey', 130);

-- Seed products.
INSERT OR IGNORE INTO products (
  name,
  icon_key,
  group_id,
  quantity_value,
  quantity_unit,
  min_quantity_value,
  missing,
  created_at,
  updated_at
)
VALUES
  ('marchewka', 'carrot', (SELECT id FROM groups WHERE name = 'warzywa'), 1.5, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('ziemniaki', 'potato', (SELECT id FROM groups WHERE name = 'warzywa'), 1.5, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mąka razowa', 'flour', (SELECT id FROM groups WHERE name = 'mąki'), 2.0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jaja', 'eggs', (SELECT id FROM groups WHERE name = 'jajka'), 30.0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mleko', 'milk', (SELECT id FROM groups WHERE name = 'nabiał'), 1.0, 'litr', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('worki na śmiecie opakowanie 35l', 'cart', (SELECT id FROM groups WHERE name = 'domowe'), 1.0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('worki na śmiecie opakowanie 60l', 'cart', (SELECT id FROM groups WHERE name = 'domowe'), 0.5, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('papier toaletowy', 'toilet-paper', (SELECT id FROM groups WHERE name = 'domowe'), 0, 'sztuk', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('wołowe mielone', 'cow', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('wołowina gulaszowa', 'cow', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('wołowina na wywar', 'cow', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('kurczak na rosół', 'chicken', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pierś z kurczaka', 'chicken', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pałki kurczaka', 'chicken', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('schab', 'pig', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mięso mielone wieprzowe', 'pig', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('karkówka', 'pig', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('boczek', 'bacon', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pierś z indyka', 'turkey', (SELECT id FROM groups WHERE name = 'mięso'), 0, 'kg', 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
