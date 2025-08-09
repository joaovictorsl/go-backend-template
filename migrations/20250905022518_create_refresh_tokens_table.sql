-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
  user_id UUID NOT NULL,
  value VARCHAR(32) NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
