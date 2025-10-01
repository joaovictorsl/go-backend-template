UPDATE refresh_tokens
SET value = $2, expires_at = $3
WHERE value = $1;
