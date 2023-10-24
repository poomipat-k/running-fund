-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
  id SERIAL PRIMARY KEY NOT NULL,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
