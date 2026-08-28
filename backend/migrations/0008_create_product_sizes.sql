-- Per-size stock, replacing the single products.stock total. Backfills from
-- the old products.sizes array, copying the old total stock to each size —
-- an approximation, since the old schema never tracked stock per size.
CREATE TABLE IF NOT EXISTS product_sizes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size TEXT NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    UNIQUE (product_id, size)
);

CREATE INDEX IF NOT EXISTS idx_product_sizes_product ON product_sizes (product_id);

INSERT INTO product_sizes (product_id, size, stock)
SELECT p.id, unnest(p.sizes), p.stock
FROM products p
WHERE array_length(p.sizes, 1) > 0;
