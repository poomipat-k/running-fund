-- +goose Up
CREATE TABLE how_to_create(
    id SERIAL PRIMARY KEY NOT NULL,
    header VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    website_config_id INT REFERENCES website_config (id)
);

-- +goose Down
ALTER TABLE how_to_create DROP COLUMN website_config_id;
DROP TABLE how_to_create;