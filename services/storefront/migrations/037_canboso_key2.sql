-- Point the canboso-backed plans at the current buyer key's catalogue.
-- Previous product ids no longer exist under either stored key.
UPDATE catalog_plans
SET supplier = 'canboso',
    supplier_product = '6ae76c35097b807a3e41cd23',
    supplier_ceiling_cents = 1105
WHERE id = 'plus';

-- The current key has no US 2FA Gmail in stock. Leave the plan unprocured
-- rather than silently shipping a Vietnam/new account against the US/2FA copy.
UPDATE catalog_plans
SET supplier = '',
    supplier_product = '',
    supplier_ceiling_cents = NULL
WHERE id = 'gmail-account';
