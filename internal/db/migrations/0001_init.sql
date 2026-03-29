-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
  id UUID PRIMARY KEY,
  service_name TEXT NOT NULL,
  price INTEGER NOT NULL CHECK (price >= 0),
  user_id UUID NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS subscriptions_user_id_idx ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS subscriptions_service_name_idx ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS subscriptions_start_date_idx ON subscriptions(start_date);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;
