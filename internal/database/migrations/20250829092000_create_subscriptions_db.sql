-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id UUID NOT NULL,
    cost FLOAT NOT NULL,
    start_date TIMESTAMP NOT NULL ,
    end_date TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
