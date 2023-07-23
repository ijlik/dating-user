-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS profiles (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name TEXT NULL,
    birth_date TIMESTAMP NULL,
    gender VARCHAR(10) NULL,
    photos TEXT NULL,
    hobby TEXT NULL,
    interest TEXT NULL,
    location TEXT NULL,
    is_premium BOOLEAN NOT NULL DEFAULT False,
    is_premium_valid_until TIMESTAMP NULL,
    daily_swap_quota INT DEFAULT 10,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS profiles;