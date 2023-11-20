-- +goose Up
CREATE TABLE review_criteria (
  id SERIAL PRIMARY KEY NOT NULL,
  criteria_version  SMALLINT NOT NULL,
  group_number SMALLINT,
  in_group_number SMALLINT,
  order_number SMALLINT NOT NULL,
  question_text VARCHAR(512),
  display_text VARCHAR(512)
);

ALTER TABLE review_criteria ADD CONSTRAINT uq_criteria_version_order_number UNIQUE(criteria_version, order_number);

-- +goose Down
DROP TABLE review_criteria;
