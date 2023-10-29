-- +goose Up
CREATE TABLE users (
  id SERIAL PRIMARY KEY NOT NULL,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255) UNIQUE NOT NULL,
  user_role VARCHAR(64)
);

-- +goose Down
DROP TABLE users;