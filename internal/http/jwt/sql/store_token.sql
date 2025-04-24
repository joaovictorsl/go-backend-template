INSERT INTO refresh_tokens (user_id, jwt, created_at, expires_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) 
DO UPDATE 
SET jwt = EXCLUDED.jwt, 
    created_at = EXCLUDED.created_at, 
    expires_at = EXCLUDED.expires_at;

