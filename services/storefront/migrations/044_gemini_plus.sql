-- Gemini goes back on sale as the supplier's 6-month Gemini AI Plus credential.
-- A new plan id is used rather than repurposing `gemini-pro-18m`: that row is a
-- claim-link product, and the supplier ships account material, so the two need
-- different delivery validation. The old row stays for its historical orders.
--
-- Ceiling 55 USD cents sits above the live 49-cent quote and ceil(55 * 6.7) =
-- 369 CNY cents stays far under the 2900 CNY price, so quote() accepts it. A
-- rise past 0.55 USD needs a new migration.
INSERT INTO catalog_plans
  (id, product_id, name, price_cents, default_cost_cents, enabled, partner_enabled,
   renewal_only, supports_warranty, delivery_kind, supplier, supplier_product,
   supplier_ceiling_cents, sort_order, stock, purchase_note, delivery_instructions)
VALUES
  ('gemini-plus-6m', 'gemini', 'Gemini Plus · 6 个月', 2900, 330, true, false,
   false, false, 'cdk', 'delivery_supplier', 'S_A74FEDC6ACCA8F00',
   55, 95, 3,
   '付款后自动发货，交付内容为 Gemini 账号资料（邮箱与密码，如含 2FA 密钥会一并提供）。',
   '用交付内容中的账号登录 gemini.google.com。请首次登录后立即修改密码，并妥善保管，不要转发给他人。');
