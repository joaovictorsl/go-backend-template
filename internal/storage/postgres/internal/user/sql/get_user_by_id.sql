SELECT id, email, created_at, updated_at
FROM users
WHERE id=$1 AND deleted_at IS NULL;
