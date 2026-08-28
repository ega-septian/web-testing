CREATE TABLE IF NOT EXISTS dress_styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_dress_styles_order ON dress_styles (display_order);

-- Same 4 categories previously hardcoded in HomeView.vue.
INSERT INTO dress_styles (name, display_order) VALUES
    ('Casual', 1),
    ('Formal', 2),
    ('Party',  3),
    ('Gym',    4);
