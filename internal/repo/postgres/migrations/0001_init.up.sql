CREATE TYPE transaction_status AS ENUM ('created', 'in_process', 'completed', 'cancelled');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    pwd_hash VARCHAR(255) NOT NULL
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status transaction_status NOT NULL
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);