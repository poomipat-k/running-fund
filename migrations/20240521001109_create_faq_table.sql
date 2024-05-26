-- +goose Up
CREATE TABLE faq(
    id SERIAL PRIMARY KEY NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(1024) NOT NULL,
    website_config_id INT REFERENCES website_config (id)
);

-- +goose Down
ALTER TABLE faq DROP COLUMN website_config_id;
DROP TABLE faq;