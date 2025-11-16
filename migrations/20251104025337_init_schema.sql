-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
                                     name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     team_name TEXT REFERENCES teams(name) ON DELETE SET NULL,
                                     is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);

CREATE TABLE IF NOT EXISTS pull_requests (
                                             id UUID PRIMARY KEY,
                                             name TEXT NOT NULL,
                                             author_id UUID NOT NULL REFERENCES users(id),
                                             created_at TIME NOT NULL DEFAULT NOW(),
                                             merged_at TIME
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);

CREATE TABLE IF NOT EXISTS pull_request_reviewers (
                                                      pull_request_id UUID NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
                                                      reviewer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                                      PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_pr_reviewers_reviewer_id;
DROP TABLE IF EXISTS pull_request_reviewers;
DROP INDEX IF EXISTS idx_pull_requests_author_id;
DROP TABLE IF EXISTS pull_requests;
DROP INDEX IF EXISTS idx_users_team_name;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd