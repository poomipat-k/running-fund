-- +goose Up
CREATE TABLE how_to_create(
    id SERIAL PRIMARY KEY NOT NULL,
    header VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    order_number INT,
    website_config_id INT REFERENCES website_config (id)
);

CREATE INDEX how_to_create_website_config_id ON how_to_create(website_config_id);

-- +goose Down
ALTER TABLE how_to_create DROP COLUMN website_config_id;
DROP TABLE how_to_create;