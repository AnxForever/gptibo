-- Several buyer keys can now be online at once: a plan names the account that
-- fulfils it, and each canboso order freezes that account so a later catalogue
-- edit can never move an in-flight purchase to a different wallet. The empty
-- account keeps meaning CANBOSO_KEY_FILE, so rows written before this migration
-- still resolve.
ALTER TABLE catalog_plans
  ADD COLUMN supplier_account text NOT NULL DEFAULT ''
  CHECK (supplier_account = '' OR supplier_account ~ '^[a-z0-9][a-z0-9_-]{0,31}$');

ALTER TABLE canboso_orders
  ADD COLUMN supplier_account text NOT NULL DEFAULT ''
  CHECK (supplier_account = '' OR supplier_account ~ '^[a-z0-9][a-z0-9_-]{0,31}$');

-- plus moves off key2's $9.94 all-warranty SKU. key1's VIP SKU is the same
-- 1-month GGPay account with a 1-month warranty at $6.34 (101 in stock when
-- checked, 5.8k sold). The ceiling leaves room for a $7.00 restock while still
-- staying far under the CNY sale price. Cheaper alternatives and their stock at
-- the same check: key4 KING $4.98 (11), key4 GGPay $4.57 (33, 24h warranty),
-- key1 GGPay $4.20 (136, 24h warranty), key2 GGpay $6.86 (151) as a standby.
UPDATE catalog_plans
SET supplier = 'canboso',
    supplier_account = 'key1',
    supplier_product = '6af5f86797bebbcda86b6cf0',
    supplier_ceiling_cents = 700
WHERE id = 'plus';
