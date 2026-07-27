-- +goose Up
CREATE TABLE watchers (
    path          TEXT PRIMARY KEY,
    first_seen_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE reported_files (
    watcher_path TEXT NOT NULL,
    file_path    TEXT NOT NULL,
    reported_at  TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (watcher_path, file_path)
);

-- +goose Down
DROP TABLE reported_files;
DROP TABLE watchers;
