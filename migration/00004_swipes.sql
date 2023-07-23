-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS swipes (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    swiper_id uuid NOT NULL,
    swiped_id uuid NOT NULL,
    is_like BOOLEAN NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (swiper_id) REFERENCES profiles (id) ON DELETE CASCADE,
    FOREIGN KEY (swiped_id) REFERENCES profiles (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS swipes;
