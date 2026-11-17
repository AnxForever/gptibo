-- Launch line: one ElevenLabs Creator 3-month credential from the delivery
-- supplier. The upstream service id is opaque; the provider identity stays out
-- of the source tree.
--
-- Cost ceiling 480 USD cents satisfies both gates at once: it sits above the
-- live 450-cent quote so the wallet check and the fulfilment charge check pass,
-- and ceil(480 * 6.7) = 3216 CNY cents stays far under the 6690 CNY price, so
-- quote() does not reject the plan as mispriced. A rise to 4.80 USD would need
-- a new migration.
--
-- The delivered value is `email:password:2fa:verify-link`, which is not a URL,
-- so it can only ship as a cdk. It stays a cdk row and the email/UI copy calls
-- it 账号资料 rather than a redemption code.
INSERT INTO catalog_plans
  (id, product_id, name, price_cents, default_cost_cents, enabled, partner_enabled,
   renewal_only, supports_warranty, delivery_kind, supplier, supplier_product,
   supplier_ceiling_cents, sort_order, purchase_note, delivery_instructions)
VALUES
  ('elevenlabs-creator-3m', 'elevenlabs', 'Creator · 3 个月', 6690, 3015, false, false,
   false, false, 'cdk', 'delivery_supplier', 'S_4BBD996F7F386644',
   480, 110,
   '付款后自动发货，交付内容为账号邮箱、密码、2FA 密钥与邮箱验证链接。',
   '四段内容以英文冒号 : 分隔，依次为邮箱、密码、2FA 密钥、邮箱验证链接。建议首次登录后立即修改密码，并妥善保管，不要转发给他人。');

-- Everything else is off while only this one line is funded. The rows stay:
-- payments/testing.go needs the `plus` row to exist for a 1-cent test order,
-- and gemini/duolingo rows still serve historical orders and the legacy
-- mappings fallback.
UPDATE catalog_plans
SET enabled = false, version = version + 1, updated_at = now()
WHERE id IN ('plus', 'pro-5x', 'pro-20x', 'pro-25x', 'codex-sms-once',
             'gemini-pro-18m', 'duolingo-12m', 'gmail-account',
             'x-premium-3m', 'x-premium-6m');
