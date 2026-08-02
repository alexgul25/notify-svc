-- +goose Up
CREATE TABLE IF NOT EXISTS inbox (
    id TEXT PRIMARY KEY,
    topic TEXT NOT NULL,
    process_status VARCHAR(50) NOT NULL,
    payload BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS inbox;