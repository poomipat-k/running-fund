-- +goose Up
CREATE TABLE website_config(
    id SERIAL PRIMARY KEY NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    landing_page TEXT,
    footer_email VARCHAR(255),
    footer_phone_number VARCHAR(64),
    footer_operate_hour VARCHAR(64)
);
-- +goose Down
DROP TABLE website_config;