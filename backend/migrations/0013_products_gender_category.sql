-- Adds gender/category/subcategory for the Shop page's filter sidebar
-- (Gender, Kategori, Sub Kategori — Ukuran already comes from product_sizes).
ALTER TABLE products
    ADD COLUMN gender TEXT NOT NULL DEFAULT '',
    ADD COLUMN category TEXT NOT NULL DEFAULT '',
    ADD COLUMN subcategory TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_products_gender ON products (gender);
CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);
CREATE INDEX IF NOT EXISTS idx_products_subcategory ON products (subcategory);
