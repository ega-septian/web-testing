-- Adds brand (e.g. "SUKO"), and drops rating entirely — there's no review
-- system behind it, so a star rating was never real data.
ALTER TABLE products
    ADD COLUMN brand TEXT NOT NULL DEFAULT '',
    DROP COLUMN rating;
