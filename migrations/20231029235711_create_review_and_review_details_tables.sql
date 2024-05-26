-- +goose Up

CREATE TABLE improvement (
  id SERIAL PRIMARY KEY NOT NULL,
  benefit BOOLEAN,
  experience_and_reliability BOOLEAN,
  fund_and_output BOOLEAN,
  project_quality BOOLEAN,
  project_standard BOOLEAN,
  vision_and_image BOOLEAN
);

CREATE TABLE review (
  id SERIAL PRIMARY KEY NOT NULL,
  user_id INT REFERENCES users(id) NOT NULL,
  project_history_id INT REFERENCES project_history(id) NOT NULL,
  is_interested_person BOOLEAN NOT NULL,
  interested_person_type VARCHAR(64),
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now(),
  summary VARCHAR(64) NOT NULL,
  improvement_id INT REFERENCES improvement(id),
  comment VARCHAR(512),
  CONSTRAINT uq_user_id_project_history_id UNIQUE(user_id, project_history_id)
);


CREATE TABLE review_details (
  id SERIAL PRIMARY KEY NOT NULL,
  review_id INT REFERENCES review(id), 
  review_criteria_id INT REFERENCES review_criteria(id),
  score SMALLINT NOT NULL
);


-- +goose Down
ALTER TABLE review DROP COLUMN user_id;
ALTER TABLE review DROP COLUMN project_history_id;
ALTER TABLE review DROP COLUMN improvement_id;

ALTER TABLE review_details DROP COLUMN review_id;
ALTER TABLE review_details DROP COLUMN review_criteria_id;

DROP TABLE improvement;
DROP TABLE review;
DROP TABLE review_details;