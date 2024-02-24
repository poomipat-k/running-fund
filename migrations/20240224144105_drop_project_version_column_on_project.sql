-- +goose Up
ALTER TABLE project DROP COLUMN project_version;

-- +goose Down
ALTER TABLE project ADD project_version SMALLINT DEFAULT 1;
