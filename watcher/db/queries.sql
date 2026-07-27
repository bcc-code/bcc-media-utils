-- name: RegisterWatcher :execrows
INSERT INTO watchers (path) VALUES (?) ON CONFLICT (path) DO NOTHING;

-- name: ListReportedFiles :many
SELECT file_path FROM reported_files WHERE watcher_path = ?;

-- name: AddReportedFile :exec
INSERT INTO reported_files (watcher_path, file_path) VALUES (?, ?)
ON CONFLICT (watcher_path, file_path) DO NOTHING;

-- name: RemoveReportedFile :exec
DELETE FROM reported_files WHERE watcher_path = ? AND file_path = ?;
