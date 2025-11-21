-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_requests (
    pull_request_id BIGSERIAL PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id BIGINT NOT NULL REFERENCES users(user_id),
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pull_requests;
-- +goose StatementEnd
