CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    section TEXT NOT NULL CHECK (section IN ('new_arrivals', 'top_selling')),
    icon TEXT NOT NULL,
    rating NUMERIC(2,1) NOT NULL,
    price INTEGER NOT NULL,
    old_price INTEGER,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_section_order ON products (section, display_order);

-- Seeds the same demo catalog that was previously hardcoded in HomeView.vue,
-- so switching the homepage to fetch from the API is a no-op visually.
INSERT INTO products (name, section, icon, rating, price, old_price, display_order) VALUES
    ('T-shirt with Tape Details', 'new_arrivals', '👕', 4.5, 120, NULL, 1),
    ('Skinny Fit Jeans',          'new_arrivals', '👖', 3.5, 240, 260,  2),
    ('Checkered Shirt',           'new_arrivals', '🧥', 4.5, 180, NULL, 3),
    ('Sleeve Striped T-shirt',    'new_arrivals', '👕', 4.5, 130, 160,  4),
    ('Vertical Striped Shirt',    'top_selling',  '🧥', 5.0, 212, 232,  1),
    ('Courage Graphic T-shirt',   'top_selling',  '👕', 4.0, 145, NULL, 2),
    ('Loose Fit Bermuda Shorts',  'top_selling',  '🩳', 3.0, 80,  NULL, 3),
    ('Faded Skinny Jeans',        'top_selling',  '👖', 4.5, 210, NULL, 4);
