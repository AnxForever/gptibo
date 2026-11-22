-- Owner repriced the launch line from 6690 to 5900 CNY cents. A ready-made
-- account that ships in one shot with no warranty and a 30-minute claim window
-- does not carry the servicing risk that backs a subscription price.
--
-- quote() still accepts the plan: the cost ceiling is ceil(480 * 6.7) = 3216
-- CNY cents, which stays under 5900. Margin is 5900 - 3015 = 2885 CNY cents.
-- Warranty stays off, so refund requests are refused for this plan; see 041.
UPDATE catalog_plans
SET price_cents = 5900, version = version + 1, updated_at = now()
WHERE id = 'elevenlabs-creator-3m' AND price_cents <> 5900;
