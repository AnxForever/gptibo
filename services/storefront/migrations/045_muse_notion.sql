-- Two more SKUs from the same supplier shelf. Both ship account material, not a
-- claim link or a redemption code, so they use delivery_kind 'cdk' and take the
-- generic account mail (see credentialPlans in internal/delivery/delivery.go).
--
-- Advertised stock is a promise the prepaid wallet can honour: every purchase is
-- deducted from the supplier balance, so the counts on sale must add up to less
-- than what is left there. The balance at the time of writing is 20.50 USD, so
-- this migration also trims Creator from 4 to 2 units. 2 * 4.50 + 3 * 0.49 +
-- 4 * 0.84 + 6 * 0.48 = 16.71 USD, leaving room for the FX buffer.
UPDATE catalog_plans
SET stock = 2, version = version + 1, updated_at = now()
WHERE id = 'elevenlabs-creator-3m' AND stock > 2;

INSERT INTO catalog_plans
  (id, product_id, name, price_cents, default_cost_cents, enabled, partner_enabled,
   renewal_only, supports_warranty, delivery_kind, supplier, supplier_product,
   supplier_ceiling_cents, sort_order, stock, purchase_note, delivery_instructions)
VALUES
  ('muse-account', 'muse', 'Muse AI · 高级账号', 3900, 563, true, false,
   false, false, 'cdk', 'delivery_supplier', 'S_E0D9F5D6FD463074',
   95, 96, 4,
   '付款后自动发货，交付内容为 Muse AI 账号资料。不支持售后质保：请务必在付款后 30 分钟内登录并修改密码，超过 30 分钟不再受理售后。',
   '用交付内容中的账号直接登录，首次登录后立即修改密码。账号为现成账号、不支持质保，请自行核对登录结果，并妥善保管资料，不要转发给他人。'),
  ('notion-business-3m', 'notion', 'Notion Business · 3 个月', 2900, 322, true, false,
   false, false, 'cdk', 'delivery_supplier', 'S_A87D8B304F46E5BD',
   60, 97, 6,
   '付款后自动发货。交付内容为开通所需资料：若为邀请链接，请用你自己的 Notion 账号打开并接受邀请；若为账号资料，直接登录即可。',
   '按交付内容操作：邀请链接请用你自己的 Notion 账号打开并接受邀请，账号资料请直接登录。开通后请自行确认工作空间权限，并妥善保管资料，不要转发给他人。');
