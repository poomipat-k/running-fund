-- +goose Up
CREATE TABLE users (
  id SERIAL PRIMARY KEY NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(128) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  user_role VARCHAR(64) DEFAULT 'applicant' NOT NULL,
  activated BOOLEAN DEFAULT false NOT NULL,
  activate_before TIMESTAMP WITH TIME ZONE DEFAULT now() + '1 day' NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
-- +goose Down
DROP TABLE users;