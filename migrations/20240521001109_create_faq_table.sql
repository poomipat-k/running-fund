-- +goose Up
CREATE TABLE faq(
    id SERIAL PRIMARY KEY NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(1024) NOT NULL,
    order_number INT,
    website_config_id INT REFERENCES website_config (id)
);

CREATE INDEX faq_website_config_id ON faq(website_config_id);

-- +goose Down
ALTER TABLE faq DROP COLUMN website_config_id;
DROP TABLE faq;