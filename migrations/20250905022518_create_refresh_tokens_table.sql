-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
  user_id UUID NOT NULL,
  value UUID NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
