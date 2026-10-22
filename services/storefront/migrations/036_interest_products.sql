-- The launch-interest list now covers the product families being sourced from
-- the reseller channels, not just the original three.
ALTER TABLE launch_interest DROP CONSTRAINT launch_interest_product_check;
ALTER TABLE launch_interest ADD CONSTRAINT launch_interest_product_check
  CHECK (product IN ('claude', 'grok', 'gmail', 'muse', 'canva', 'cursor', 'notion', 'perplexity', 'manus', 'framer'));
