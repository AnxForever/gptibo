-- The upstream listing sells a ready-made, stable account with no warranty and a
-- 30-minute claim window: the buyer must sign in and change the password within
-- 30 minutes of purchase, and only claims filed inside that window are accepted.
-- Neither the stability nor the two rules were in the 040 copy, so they are
-- restated here where the buyer actually reads them.
UPDATE catalog_plans
SET purchase_note = '现成稳定账号，付款确认后自动发货。不支持售后质保：请务必在付款后 30 分钟内登录并修改账号密码，超过 30 分钟不再受理售后。',
    delivery_instructions = '四段内容依次为邮箱、邮箱密码、Creator 密码、邮箱查看链接，之间以英文冒号分隔；最后一段是完整的查看链接，从 https:// 开始，请整段复制。
请在付款后 30 分钟内登录并修改账号密码——只有这 30 分钟内提交的售后才受理。账号为现成稳定账号（131K 积分/月），不支持质保，请自行核对登录结果，并妥善保管资料，不要转发给他人。',
    version = version + 1,
    updated_at = now()
WHERE id = 'elevenlabs-creator-3m';
