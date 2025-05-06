-- +goose Up
CREATE TABLE operation_config(
    id SERIAL PRIMARY KEY NOT NULL,
    allow_new_project BOOLEAN NOT NULL
);

INSERT INTO operation_config (allow_new_project)
VALUES (FALSE);

-- +goose Down
DROP TABLE operation_config;