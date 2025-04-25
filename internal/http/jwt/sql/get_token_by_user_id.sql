SELECT user_id, jwt, created_at, expires_at FROM refresh_tokens WHERE user_id=$1;
