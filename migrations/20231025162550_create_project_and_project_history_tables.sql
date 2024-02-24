-- +goose Up
-- TABLES
CREATE TABLE project (
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255) UNIQUE NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT now(),
  project_history_id INT,
  user_id INT
  
);

CREATE TABLE project_history(
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255)  NOT NULL,
  project_version  SMALLINT DEFAULT 1,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  download_link VARCHAR(512),
  admin_comment VARCHAR(512),
  project_id INT,
  -- STEP 0
  collaborated INT NOT NULL,
  -- STEP 1
  project_name VARCHAR(512) NOT NULL,
  from_date TIMESTAMP WITH TIME ZONE NOT NULL,
  to_date TIMESTAMP WITH TIME ZONE NOT NULL,
  address_id INT,
  start_point VARCHAR(255) NOT NULL,
  finish_point VARCHAR(255) NOT NULL,
  cat_road_race BOOLEAN NOT NULL,
  cat_trail_running BOOLEAN NOT NULL,
  cat_has_other BOOLEAN NOT NULL,
  cat_other_type VARCHAR(255) NOT NULL,
  vip BOOLEAN NOT NULL,
  expected_participants INT NOT NULL,
  has_organizer BOOLEAN NOT NULL,
  organizer_name VARCHAR(255) NOT NULL
);

CREATE TABLE distance(
  id SERIAL PRIMARY KEY NOT NULL,
  type VARCHAR(255) NOT NULL,
  fee FLOAT NOT NULL,
  is_dynamic BOOLEAN NOT NULL,
  project_history_id INT
);

-- CONSTRAINTS
ALTER TABLE project ADD CONSTRAINT fk_project_history_project FOREIGN KEY (project_history_id) REFERENCES project_history (id);

ALTER TABLE project ADD CONSTRAINT fk_users_project FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE project_history ADD CONSTRAINT fk_project_project_history FOREIGN KEY (project_id) REFERENCES project(id);

ALTER TABLE project_history ADD CONSTRAINT fk_address_project_history FOREIGN KEY (address_id) REFERENCES address (id);

ALTER TABLE distance ADD CONSTRAINT fk_project_history_distance FOREIGN KEY (project_history_id) REFERENCES project_history (id);

CREATE INDEX project_created_at ON project (created_at);

-- +goose Down
ALTER TABLE project DROP COLUMN project_history_id;
ALTER TABLE project DROP COLUMN user_id;

ALTER TABLE project_history DROP COLUMN address_id;
ALTER TABLE project_history DROP COLUMN project_id;

ALTER TABLE distance DROP COLUMN project_history_id;

DROP TABLE project;
DROP TABLE project_history;
DROP TABLE distance;
