-- +goose Up
CREATE TABLE review_period (
    id SERIAL PRIMARY KEY NOT NULL,
    from_date TIMESTAMP WITH TIME ZONE NOT NULL,
    to_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE review_period;