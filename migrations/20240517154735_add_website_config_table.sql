-- +goose Up
CREATE TABLE website_config(
    id SERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    landing_page TEXT
);
-- +goose Down
DROP TABLE website_config;