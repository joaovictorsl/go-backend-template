SELECT u.id, u.email, u.created_at, u.updated_at
FROM linked_accounts la
JOIN users u ON la.user_id = u.id
WHERE la.provider=$1 AND la.provider_user_id=$2;
