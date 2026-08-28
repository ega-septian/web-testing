-- Products are now represented by their uploaded photos (product_images),
-- not a manually chosen emoji.
ALTER TABLE products DROP COLUMN icon;
