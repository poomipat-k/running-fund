-- +goose Up
CREATE TABLE users (
  id SERIAL PRIMARY KEY NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  user_role VARCHAR(64) DEFAULT 'user' NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now() NOT NULL
);

-- +goose Down
DROP TABLE users;