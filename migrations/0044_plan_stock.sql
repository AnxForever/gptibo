-- Mirrors the Go backend: catalog_plans.stock is the sellable count a shop card
-- advertises. NULL hides the count. The worker re-reads it with SELECT *, so the
-- catalog response carries it without further changes.
ALTER TABLE catalog_plans ADD COLUMN stock INTEGER;
