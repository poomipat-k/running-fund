-- +goose Up
-- TABLES
CREATE TABLE project_history(
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255)  NOT NULL,
  project_version  SMALLINT DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  download_link VARCHAR(512),
  admin_comment VARCHAR(512),
  -- STEP 0
  collaborated BOOLEAN NOT NULL,
  -- STEP 1
  project_name VARCHAR(512) NOT NULL,
  from_date TIMESTAMP WITH TIME ZONE NOT NULL,
  to_date TIMESTAMP WITH TIME ZONE NOT NULL,
  address_id INT, CONSTRAINT fk_project_history_address FOREIGN KEY (address_id) REFERENCES address (id),
  start_point VARCHAR(255) NOT NULL,
  finish_point VARCHAR(255) NOT NULL,
  cat_road_race BOOLEAN NOT NULL,
  cat_trail_running BOOLEAN NOT NULL,
  cat_has_other BOOLEAN NOT NULL,
  cat_other_type VARCHAR(255) NOT NULL,
  vip BOOLEAN NOT NULL,
  expected_participants VARCHAR(64) NOT NULL,
  has_organizer BOOLEAN NOT NULL,
  organizer_name VARCHAR(255) NOT NULL,
  -- STEP 2
  project_head_contact_id INT, CONSTRAINT fk_project_history_contact_project_head FOREIGN KEY (project_head_contact_id) REFERENCES contact(id),
  project_manager_contact_id INT, CONSTRAINT fk_project_history_contact_project_manager FOREIGN KEY (project_manager_contact_id) REFERENCES contact(id),
  project_coordinator_contact_id INT, CONSTRAINT fk_project_history_contact_project_coordinator FOREIGN KEY (project_coordinator_contact_id) REFERENCES contact(id),
  project_race_director_contact_id INT, CONSTRAINT fk_project_history_contact_project_race_director FOREIGN KEY (project_race_director_contact_id) REFERENCES contact(id),
  organization_type VARCHAR(255) NOT NULL,
  organization_name VARCHAR(255) NOT NULL,
  -- STEP 3
  background TEXT NOT NULL,
  objective TEXT NOT NULL,
  mkt_has_facebook BOOLEAN NOT NULL,
  mkt_facebook VARCHAR(255),
  mkt_has_website BOOLEAN NOT NULL,
  mkt_website VARCHAR(255),
  mkt_use_online_page BOOLEAN NOT NULL,
  mkt_online_page VARCHAR(255),
  mkt_use_other_online_marketing BOOLEAN NOT NULL,
  mkt_other_online_marketing VARCHAR(255),
  mkt_pr BOOLEAN NOT NULL,
  mkt_local_official BOOLEAN NOT NULL,
  mkt_booth BOOLEAN NOT NULL,
  mkt_billboard BOOLEAN NOT NULL,
  mkt_tv BOOLEAN NOT NULL,
  mkt_use_other_offline_marketing BOOLEAN NOT NULL,
  mkt_other_offline_marketing VARCHAR(255),
  st_runner_info BOOLEAN NOT NULL,
  st_health_decider BOOLEAN NOT NULL,
  st_ambulance BOOLEAN NOT NULL,
  st_first_aid BOOLEAN NOT NULL,
  st_aed BOOLEAN NOT NULL,
  st_aed_count INT,
  st_insurance BOOLEAN NOT NULL,
  st_other BOOLEAN NOT NULL,
  st_addition VARCHAR(255),
  measure_athletics_association BOOLEAN NOT NULL,
  measure_calibrated_bicycle BOOLEAN NOT NULL,
  measure_self_measurement BOOLEAN NOT NULL,
  measure_self_tool VARCHAR(255),
  traffic_ask_permission BOOLEAN NOT NULL,
  traffic_has_supporter BOOLEAN NOT NULL,
  traffic_road_closure BOOLEAN NOT NULL,
  traffic_signs BOOLEAN NOT NULL,
  traffic_lighting BOOLEAN NOT NULL,
  judge_type VARCHAR(255) NOT NULL,
  judge_other_type VARCHAR(255),
  support_provincial_admin BOOLEAN NOT NULL,
  support_safety BOOLEAN NOT NULL,
  support_health BOOLEAN NOT NULL,
  support_volunteer BOOLEAN NOT NULL,
  support_community BOOLEAN NOT NULL,
  support_other BOOLEAN NOT NULL,
  support_addition VARCHAR(255),
  feedback TEXT NOT NULL,
  -- STEP 4
  exp_this_first_time BOOLEAN NOT NULL,
  exp_this_ordinal_number INT NOT NULL,
  exp_this_latest_date TIMESTAMP WITH TIME ZONE NOT NULL,
  exp_this_completed1_year SMALLINT NOT NULL,
  exp_this_completed1_name VARCHAR(255) NOT NULL,
  exp_this_completed1_participant INT NOT NULL,
  exp_this_completed2_year SMALLINT,
  exp_this_completed2_name VARCHAR(255),
  exp_this_completed2_participant INT,
  exp_this_completed3_year SMALLINT,
  exp_this_completed3_name VARCHAR(255),
  exp_this_completed3_participant INT,
  exp_other_done_before BOOLEAN NOT NULL,
  exp_other_completed1_year SMALLINT NOT NULL,
  exp_other_completed1_name VARCHAR(255) NOT NULL,
  exp_other_completed1_participant INT NOT NULL,
  exp_other_completed2_year SMALLINT,
  exp_other_completed2_name VARCHAR(255),
  exp_other_completed2_participant INT,
  exp_other_completed3_year SMALLINT,
  exp_other_completed3_name VARCHAR(255),
  exp_other_completed3_participant INT,
  -- STEP 5
  fund_total INT NOT NULL,
  fund_support_organization TEXT NOT NULL,
  fund_req_fund BOOLEAN NOT NULL,
  fund_req_fund_amount INT,
  fund_req_bib BOOLEAN NOT NULL,
  fund_req_bib_amount INT,
  fund_req_pr BOOLEAN NOT NULL,
  fund_req_seminar BOOLEAN NOT NULL,
  fund_req_seminar_topic VARCHAR(255),
  fund_req_other BOOLEAN NOT NULL,
  fund_req_other_type VARCHAR(255),
  files_prefix VARCHAR(255) NOT NULL
);

CREATE TABLE project (
  id SERIAL PRIMARY KEY NOT NULL,
  project_code VARCHAR(255) UNIQUE NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL  DEFAULT now(),
  project_history_id INT ,CONSTRAINT fk_project_project_history FOREIGN KEY(project_history_id) REFERENCES project_history (id),
  user_id INT ,CONSTRAINT fk_project_users FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE distance(
  id SERIAL PRIMARY KEY NOT NULL,
  type VARCHAR(255) NOT NULL,
  fee FLOAT NOT NULL,
  is_dynamic BOOLEAN NOT NULL,
  project_history_id INT,
  CONSTRAINT fk_distance_project_history FOREIGN KEY (project_history_id) REFERENCES project_history (id)
);

-- INDEX
CREATE INDEX project_created_at ON project (created_at);

-- +goose Down
ALTER TABLE project DROP COLUMN project_history_id;
ALTER TABLE project DROP COLUMN user_id;

ALTER TABLE project_history DROP COLUMN address_id;
ALTER TABLE project_history DROP COLUMN project_head_contact_id;
ALTER TABLE project_history DROP COLUMN project_manager_contact_id;
ALTER TABLE project_history DROP COLUMN project_coordinator_contact_id;
ALTER TABLE project_history DROP COLUMN project_race_director_contact_id;

ALTER TABLE distance DROP COLUMN project_history_id;

DROP TABLE project;
DROP TABLE project_history;
DROP TABLE distance;
