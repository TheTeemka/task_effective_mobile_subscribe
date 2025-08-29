-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_id ON subscriptions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_id;
-- +goose StatementEnd
