-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS one_time_password_logs (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    onetime_password_type VARCHAR(10) NOT NULL DEFAULT 'PHONE', -- "EMAIL, PHONE"
    code TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'UNUSED', -- "EXPIRED, USED, UNUSED"
    otp_limit INT DEFAULT 3,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_code_otp ON one_time_password_logs(onetime_password_type, user_id, code);

-- +goose Down
DROP INDEX IF EXISTS idx_code_otp;
DROP TABLE IF EXISTS one_time_password_logs;
