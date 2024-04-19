-- +goose Up
ALTER TABLE project_history
ADD st_volunteer_doctor BOOLEAN DEFAULT false;
-- +goose Down
ALTER TABLE project_history DROP COLUMN st_volunteer_doctor;