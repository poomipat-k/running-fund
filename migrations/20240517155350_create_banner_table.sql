-- +goose Up
CREATE TABLE banner(
    id SERIAL PRIMARY KEY NOT NULL,
    full_path VARCHAR(255) NOT NULL,
    link_to VARCHAR(512) NOT NULL,
    website_config_id INT REFERENCES website_config (id)
);

-- +goose Down
ALTER TABLE banner DROP COLUMN website_config_id;
DROP TABLE banner;