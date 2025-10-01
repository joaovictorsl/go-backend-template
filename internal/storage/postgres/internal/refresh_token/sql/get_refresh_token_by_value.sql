SELECT user_id, value, expires_at
FROM refresh_tokens
WHERE value=$1;
