-- The homepage's "Browse by Dress Style" section (Casual/Formal/Party/Gym)
-- was disconnected from the real product catalog — no product's
-- category/subcategory ever matched these values, so the cards weren't
-- actually clickable to anything. Replaced by a "Shop by Category" section
-- driven by products.category (see HomeView.vue), so this table is unused.
DROP TABLE IF EXISTS dress_styles;
