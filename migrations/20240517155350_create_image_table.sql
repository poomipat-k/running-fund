-- +goose Up
CREATE TABLE website_image(
    id SERIAL PRIMARY KEY NOT NULL,
    code VARCHAR(64) NOT NULL,
    full_path VARCHAR(255) NOT NULL,
    object_key VARCHAR(255) NOT NULL,
    link_to VARCHAR(512),
    order_number INT,
    website_config_id INT REFERENCES website_config (id)
);

CREATE INDEX website_image_website_config_id ON website_image(website_config_id);

-- +goose Down
ALTER TABLE website_image DROP COLUMN website_config_id;
DROP TABLE website_image;