-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payments (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    amount FLOAT NOT NULL,
    identifier TEXT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_data TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS payments;
