SELECT user_id, jwt, create_at, expires_at FROM refresh_tokens WHERE user_id=$1;
