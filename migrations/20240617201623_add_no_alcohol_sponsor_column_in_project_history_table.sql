-- +goose Up
ALTER TABLE project_history ADD no_alcohol_sponsor BOOLEAN DEFAULT TRUE;

-- +goose Down
ALTER TABLE project_history DROP COLUMN no_alcohol_sponsor;