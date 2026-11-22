-- Owner-set pricing: the three low-cost lines earn a flat 10 CNY each, while
-- Creator keeps its 59 CNY price because its upstream cost is 30.15 CNY.
--
--   gemini-plus-6m      3.28 cost -> 13.90 (10.62 margin)
--   muse-account        5.63 cost -> 15.90 (10.27 margin)
--   notion-business-3m  3.22 cost -> 13.90 (10.68 margin)
--
-- Every ceiling still converts well below its new price, so quote() keeps
-- accepting the plans: 55c -> 3.69, 95c -> 6.37, 60c -> 4.02.
UPDATE catalog_plans SET price_cents = 1390, version = version + 1, updated_at = now()
WHERE id = 'gemini-plus-6m' AND price_cents <> 1390;
UPDATE catalog_plans SET price_cents = 1590, version = version + 1, updated_at = now()
WHERE id = 'muse-account' AND price_cents <> 1590;
UPDATE catalog_plans SET price_cents = 1390, version = version + 1, updated_at = now()
WHERE id = 'notion-business-3m' AND price_cents <> 1390;
