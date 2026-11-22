-- The delivery supplier ships `email:mailbox-password:account-password:verify-link`.
-- There is no 2FA secret in the credential, so the buyer-facing copy must not
-- promise one. 039 is checksum-locked once applied; the two copy columns are
-- rewritten here instead.
UPDATE catalog_plans
SET purchase_note = '付款后自动发货，交付内容为账号邮箱、邮箱密码、账号密码与邮箱验证链接。',
    delivery_instructions = '四段内容依次为邮箱、邮箱密码、账号密码、邮箱验证链接，之间以英文冒号分隔；最后一段是完整的验证链接，从 https:// 开始，请整段复制。建议首次登录后立即修改密码，并妥善保管，不要转发给他人。',
    version = version + 1,
    updated_at = now()
WHERE id = 'elevenlabs-creator-3m';
