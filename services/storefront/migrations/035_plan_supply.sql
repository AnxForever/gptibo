-- Make the catalogue data-driven: how a plan is delivered and who fulfils it
-- live on the plan row instead of in Go switch statements.
ALTER TABLE catalog_plans
  ADD COLUMN delivery_kind text NOT NULL DEFAULT '' CHECK (delivery_kind IN ('', 'cdk', 'link', 'account')),
  ADD COLUMN supplier text NOT NULL DEFAULT '' CHECK (supplier IN ('', 'delivery_supplier', 'canboso', 'gmail')),
  ADD COLUMN supplier_product text NOT NULL DEFAULT '',
  ADD COLUMN supplier_ceiling_cents bigint CHECK (supplier_ceiling_cents IS NULL OR supplier_ceiling_cents > 0);

-- Orders snapshot both decisions so a later plan edit can never redirect an
-- existing order to a different procurement table or delivery format.
ALTER TABLE orders
  ADD COLUMN delivery_kind text NOT NULL DEFAULT '',
  ADD COLUMN procurement_kind text NOT NULL DEFAULT '';

UPDATE catalog_plans SET delivery_kind = 'link' WHERE id = 'gemini-pro-18m';
UPDATE catalog_plans SET delivery_kind = 'cdk' WHERE id IN ('duolingo-12m', 'gmail-account', 'plus');

UPDATE catalog_plans SET supplier = 'delivery_supplier', supplier_product = 'S_09', supplier_ceiling_cents = 80 WHERE id = 'gemini-pro-18m';
UPDATE catalog_plans SET supplier = 'delivery_supplier', supplier_product = 'S_03', supplier_ceiling_cents = 60 WHERE id = 'duolingo-12m';
UPDATE catalog_plans SET supplier = 'canboso', supplier_product = '6ad8a9470c1a1b0e66502c1c', supplier_ceiling_cents = 987 WHERE id = 'plus';
UPDATE catalog_plans SET supplier = 'canboso', supplier_product = '6ad90b1b0c1a1b0e66654295', supplier_ceiling_cents = 96 WHERE id = 'gmail-account';

-- Existing live orders keep the format they were created under.
UPDATE orders o SET delivery_kind = p.delivery_kind FROM catalog_plans p WHERE p.id = o.plan AND o.delivery_kind = '';
