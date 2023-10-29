-- +goose Up
CREATE TABLE project (
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255) UNIQUE NOT NULL,
  project_version  SMALLINT DEFAULT 1,
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now(),
  project_history_id INT
);

CREATE TABLE project_history(
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255)  NOT NULL,
  project_name Text NOT NULL,
  project_version  SMALLINT DEFAULT 1,
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now(),
  download_link TEXT,
  admin_comment TEXT,
  project_id INT
);

ALTER TABLE project ADD CONSTRAINT fk_project_history_project FOREIGN KEY (project_history_id) REFERENCES project_history (id);

ALTER TABLE project_history ADD CONSTRAINT fk_project_project_history FOREIGN KEY (project_id) REFERENCES project(id);

-- +goose Down
ALTER TABLE project DROP COLUMN project_history_id;
ALTER TABLE project_history DROP COLUMN project_id;

DROP TABLE project;
DROP TABLE project_history;
