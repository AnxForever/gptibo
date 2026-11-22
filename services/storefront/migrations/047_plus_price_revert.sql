-- Plus stays off sale: the only 1-month full-warranty supply sits behind an
-- unfunded canboso wallet, and the owner does not want a no-warranty listing.
-- Put the price back to its paused value so the shop advertises nothing new.
--
-- The supplier_product remap stays. The previously configured SKU was delisted
-- by the seller, so restoring it would only bring back a dead purchase path;
-- whoever reopens this plan next should start from the live row.
UPDATE catalog_plans
SET price_cents = 13900,
    version = version + 1,
    updated_at = now()
WHERE id = 'plus' AND enabled = false;
