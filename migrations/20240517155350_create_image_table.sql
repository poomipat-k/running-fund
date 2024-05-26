-- +goose Up
CREATE TABLE website_image(
    id SERIAL PRIMARY KEY NOT NULL,
    code VARCHAR(64) NOT NULL,
    full_path VARCHAR(255) NOT NULL,
    object_key VARCHAR(255) NOT NULL,
    link_to VARCHAR(512),
    website_config_id INT REFERENCES website_config (id)
);

-- +goose Down
ALTER TABLE website_image DROP COLUMN website_config_id;
DROP TABLE website_image;