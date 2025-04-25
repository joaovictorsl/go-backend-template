CREATE TABLE refresh_tokens (
    user_id INT PRIMARY KEY,
    jwt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL CHECK ( created_at <= expires_at ),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
