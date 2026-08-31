-- Soft-delete marker for products. Used to clean up test-seeded catalog
-- data (see api/products.ts's seedProduct in gasntin-automation) without
-- permanently destroying rows or breaking historical order_items that
-- reference a product (order_items.product_id is ON DELETE SET NULL, which
-- would otherwise silently orphan a real customer's order history).
ALTER TABLE products ADD COLUMN deleted_at TIMESTAMPTZ NULL;

-- Every product listing/detail/filter query filters on this being NULL, so
-- an index keeps that filter cheap as the catalog (and its deleted rows) grow.
CREATE INDEX idx_products_deleted_at ON products (deleted_at) WHERE deleted_at IS NULL;
