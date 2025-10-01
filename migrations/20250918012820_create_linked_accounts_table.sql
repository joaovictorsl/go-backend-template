-- +goose Up
-- +goose StatementBegin
CREATE TABLE linked_accounts (
  user_id UUID NOT NULL,
  provider VARCHAR(20) NOT NULL,
  provider_user_id VARCHAR(64) NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  UNIQUE (provider, provider_user_id),
  PRIMARY KEY (user_id, provider),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE linked_accounts;
-- +goose StatementEnd
