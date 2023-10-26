-- +goose Up
-- +goose StatementBegin
CREATE TABLE project (
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255) UNIQUE NOT NULL,
  project_version  SMALLINT DEFAULT 1,
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project;
-- +goose StatementEnd
