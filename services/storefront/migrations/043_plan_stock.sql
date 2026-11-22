-- Stock lets a plan advertise how many units are still sellable on the shop
-- card. NULL means the plan does not show a quantity at all, which is the
-- right state for plans whose supply we cannot count.
--
-- The procurement path decrements it once per verified upstream purchase, and
-- the workspace resets it when stock is restocked. It never gates checkout:
-- the upstream quote is still the authority on whether a unit exists.
ALTER TABLE catalog_plans ADD COLUMN stock integer CHECK (stock >= 0);
