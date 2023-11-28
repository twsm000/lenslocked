CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL,
    token BYTEA UNIQUE NOT NULL CHECK(octet_length(token) = 64),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
