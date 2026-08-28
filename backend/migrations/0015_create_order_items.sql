-- Line items for an order. product_id is nullable (ON DELETE SET NULL) so
-- deleting a product later doesn't erase order history — everything needed
-- to display the order (name, brand, price at the time of purchase) is
-- snapshotted here rather than looked up live from products, which can
-- change (or disappear) after the fact.
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    product_name TEXT NOT NULL,
    brand TEXT NOT NULL,
    size TEXT NOT NULL,
    unit_price INTEGER NOT NULL CHECK (unit_price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items (order_id);
