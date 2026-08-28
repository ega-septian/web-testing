-- `section` is retired: New Arrivals and Top Selling are now computed
-- (created_at / total_sold from the sales table), not manually tagged.
-- `sizes` and `stock` moved to product_sizes (see migration 0008).
DROP INDEX IF EXISTS idx_products_section_created;

ALTER TABLE products
    DROP COLUMN section,
    DROP COLUMN sizes,
    DROP COLUMN stock;
