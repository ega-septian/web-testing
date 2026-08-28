-- Minimal sale-event log: no orders/cart/checkout system exists yet, so this
-- is intentionally just "N units of product X sold at time T" rather than a
-- full order model. Enough to compute total_sold per product for ranking.
CREATE TABLE IF NOT EXISTS sales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    sold_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sales_product ON sales (product_id);
