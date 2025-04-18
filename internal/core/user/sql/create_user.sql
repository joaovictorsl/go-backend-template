INSERT INTO users (google_id, email, username) VALUES ($1, $2, $3) RETURNING id;
