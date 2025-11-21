-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_request_reviewers (
    pull_request_id BIGINT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pull_request_reviewers;
-- +goose StatementEnd
