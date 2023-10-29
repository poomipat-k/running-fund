-- +goose Up
CREATE TABLE review_criteria (
  id SERIAL PRIMARY KEY NOT NULL,
  criteria_version  SMALLINT,
  group_number SMALLINT,
  in_group_number SMALLINT,
  order_number SMALLINT,
  question_text TEXT
);

-- +goose Down
DROP TABLE review_criteria;
