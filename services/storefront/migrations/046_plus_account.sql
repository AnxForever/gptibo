-- ChatGPT Plus as a ready-made account at 70 CNY.
--
-- The old mapping pointed at a key1 SKU the seller has since delisted, so the
-- purchase path was dead even though the row looked configured. The live
-- equivalent is the same seller's GGpay line, still sold as 保修 1 个月
-- (0.0671 USD per unit, 47 in stock at the time of writing). The ceiling moves
-- to 0.0850 USD so a modest list-price bump does not silently stop sales while
-- staying far under the 7000 CNY price.
--
-- enabled stays false on purpose: every canboso wallet on this deployment reads
-- 0.00 USD, and a paid order that the wallet cannot settle is worse than an
-- unlisted product. Flip enabled and set stock once the wallet is funded.
UPDATE catalog_plans
SET price_cents = 7000,
    supplier_product = '6afc431ed8b8b2c2b4ba7c80',
    supplier_account = 'key1',
    supplier_ceiling_cents = 850,
    stock = NULL,
    version = version + 1,
    updated_at = now()
WHERE id = 'plus';
