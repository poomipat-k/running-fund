-- +goose Up
CREATE TABLE contact (
  id SERIAL PRIMARY KEY NOT NULL,
  prefix VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  organization_position VARCHAR(255),
  event_position VARCHAR(255),
  address_id INT, CONSTRAINT fk_contact_address FOREIGN KEY (address_id) REFERENCES address (id),
  email VARCHAR(255),
  line_id VARCHAR(255),
  phone_number VARCHAR(64)
);

-- +goose Down
ALTER TABLE contact DROP COLUMN address_id;
DROP TABLE contact;