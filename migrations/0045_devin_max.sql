-- Mirrors the Go backend: one staff-delivered Devin Max claim link at 240 CNY.
-- The worker's delivery format comes from shared/email-products.ts, so this row
-- only needs the listing itself. Stock stays NULL because the count is not
-- advertised; take the plan off sale once the single link is sold.
INSERT INTO catalog_plans(id,product_id,name,price_cents,enabled,renewal_only,updated_at,partner_enabled,purchase_note,delivery_instructions,stock)
VALUES ('devin-max-1m','devin','Devin Max · 1 个月',24000,1,0,0,0,
'仅 1 份，售出后下架。付款确认后由工作人员发送专属开通链接，可在订单页查看，也会发送至下单邮箱。链接仅限本订单使用，请勿转发。',
'用你的 Devin 账号打开领取链接，按页面提示完成开通。开通后请核对 Devin Max 订阅状态；链接仅限本订单使用，请勿转发给他人。',NULL);
