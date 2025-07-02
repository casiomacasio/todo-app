DROP TABLE IF EXISTS refresh_tokens;

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    hashed_token TEXT NOT NULL,
    issued_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);