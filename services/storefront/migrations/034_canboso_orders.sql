-- Canboso buyer-key orders: one row per local order that buys from the
-- canboso.com Telegram buyer API. The purchase response is the only place the
-- credentials appear, so whatever it returned is frozen in `delivery`.
CREATE TABLE canboso_orders (
  order_id TEXT PRIMARY KEY,
  product_id TEXT NOT NULL,
  credential_hash TEXT NOT NULL,
  request_key TEXT NOT NULL UNIQUE,
  request_body TEXT NOT NULL,
  max_usd_cents bigint NOT NULL CHECK(max_usd_cents>0),
  cny_per_usd REAL NOT NULL CHECK(cny_per_usd>0),
  instructions TEXT NOT NULL,
  state TEXT NOT NULL DEFAULT 'prepared' CHECK(state IN ('prepared','creating','unknown','pending','purchased','rejected')),
  supplier_order TEXT UNIQUE,
  actual_usd_cents bigint,
  actual_usd_micros bigint,
  actual_cost_cents bigint,
  delivery TEXT,
  value TEXT,
  attempted_at bigint,
  lease_token TEXT,
  lease_until bigint NOT NULL DEFAULT 0,
  next_check_at bigint NOT NULL DEFAULT 0,
  updated_at bigint NOT NULL
);

CREATE INDEX canboso_reconcile ON canboso_orders(state,next_check_at,updated_at);

ALTER TABLE canboso_orders ADD FOREIGN KEY (order_id) REFERENCES orders(id) DEFERRABLE INITIALLY IMMEDIATE;
