-- Devin Max as a single staff-delivered claim link at 240 CNY.
--
-- The seller holds one redemption link and pastes it into the order after
-- payment; this SKU has no upstream API. delivery_kind 'link' makes the
-- workspace validate an HTTPS claim URL and the generic delivery-link mail
-- carry it. supplier stays empty, which procurement now reads as "fulfilled by
-- staff": the plan is sellable whenever the email channel is on, but no
-- purchase row is written and nothing is bought upstream.
--
-- stock stays NULL on purpose — the seller does not want an advertised count.
-- Once the single link is sold, take the plan off sale in the workspace.
INSERT INTO catalog_plans
  (id, product_id, name, price_cents, enabled, partner_enabled,
   renewal_only, supports_warranty, delivery_kind, supplier, supplier_product,
   supplier_ceiling_cents, sort_order, stock, purchase_note, delivery_instructions)
VALUES
  ('devin-max-1m', 'devin', 'Devin Max · 1 个月', 24000, true, false,
   false, false, 'link', '', '', NULL, 98, NULL,
   '仅 1 份，售出后下架。付款确认后由工作人员发送专属开通链接，可在订单页查看，也会发送至下单邮箱。链接仅限本订单使用，请勿转发。',
   '用你的 Devin 账号打开领取链接，按页面提示完成开通。开通后请核对 Devin Max 订阅状态；链接仅限本订单使用，请勿转发给他人。');
