CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
