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
  user_id INT REFERENCES users(id),
  project_history_id INT REFERENCES project_history(id),
  is_interested_person BOOLEAN NOT NULL,
  interested_person_type VARCHAR(64),
  created_at  TIMESTAMP WITH TIME ZONE  DEFAULT now(),
  summary VARCHAR(64) NOT NULL,
  improvement_id INT REFERENCES improvement(id),
  comment VARCHAR(512)
);


CREATE TABLE review_details (
  id SERIAL PRIMARY KEY NOT NULL,
  review_id INT,
  review_criteria_id INT REFERENCES review_criteria(id),
  score SMALLINT NOT NULL
);

ALTER TABLE review_details ADD CONSTRAINT fk_review_review_details FOREIGN KEY (review_id) REFERENCES review (id);

-- +goose Down
ALTER TABLE review DROP COLUMN user_id;
ALTER TABLE review DROP COLUMN project_history_id;
ALTER TABLE review DROP COLUMN improvement_id;

ALTER TABLE review_details DROP COLUMN review_id;
ALTER TABLE review_details DROP COLUMN review_criteria_id;

DROP TABLE improvement;
DROP TABLE review;
DROP TABLE review_details;