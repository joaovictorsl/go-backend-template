INSERT INTO users (provider_id, email) VALUES ($1, $2) RETURNING id;
