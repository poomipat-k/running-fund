-- +goose Up
CREATE TABLE review_period (
    id SERIAL PRIMARY KEY NOT NULL,
    from_date  TIMESTAMP WITH TIME ZONE  NOT NULL,
    to_date  TIMESTAMP WITH TIME ZONE  NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT now()
);

INSERT INTO review_period (from_date, to_date)
VALUES ('2023-11-01 17:00:00.000000+00', '2023-12-31 17:00:00.000000+00');

-- +goose Down
DROP TABLE review_period;
