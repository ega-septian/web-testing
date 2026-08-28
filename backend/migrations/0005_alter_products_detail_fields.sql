-- Adds product-detail-page fields, and switches New Arrivals / Top Selling
-- ordering from a manual sequence to newest-first (display_order is dropped).
ALTER TABLE products
    ADD COLUMN description TEXT NOT NULL DEFAULT '',
    ADD COLUMN sizes TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN stock INTEGER NOT NULL DEFAULT 0;

ALTER TABLE products DROP COLUMN display_order;

CREATE INDEX IF NOT EXISTS idx_products_section_created ON products (section, created_at DESC);

-- Backfill the existing demo catalog with detail-page content.
UPDATE products SET description = 'Produk pilihan dari koleksi ' ||
    CASE section WHEN 'new_arrivals' THEN 'terbaru' ELSE 'terlaris' END ||
    ' SHOP.CO. Dibuat dari bahan berkualitas untuk kenyamanan sehari-hari.',
    sizes = ARRAY['S', 'M', 'L', 'XL'],
    stock = 25
WHERE description = '';
