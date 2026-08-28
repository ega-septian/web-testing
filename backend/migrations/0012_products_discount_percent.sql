-- Replaces old_price (a second absolute price to keep in sync manually)
-- with discount as a percentage — price stays the one source of truth for
-- what's actually charged; the struck-through "original price" shown in the
-- UI is derived from price + discount, not stored separately.
ALTER TABLE products
    DROP COLUMN old_price,
    ADD COLUMN discount INTEGER NOT NULL DEFAULT 0 CHECK (discount >= 0 AND discount <= 100);
