CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token BYTEA UNIQUE NOT NULL CHECK(octet_length(token) = 64)
);
