package main

import (
	"context"

	"github.com/bcc-code/bccm-utils/watcher/db"
)

// reportedStore persists the reported-files set of a single watcher so it
// survives restarts. The in-memory filesReported map stays authoritative
// during the process lifetime; the store is write-through only.
type reportedStore interface {
	Add(ctx context.Context, file string) error
	Remove(ctx context.Context, file string) error
}

type sqliteStore struct {
	q           *db.Queries
	watcherPath string
}

func newStore(q *db.Queries, watcherPath string) *sqliteStore {
	return &sqliteStore{q: q, watcherPath: watcherPath}
}

// Register records the watcher pattern in the database and reports whether
// this is the first run for it.
func (s *sqliteStore) Register(ctx context.Context) (firstRun bool, err error) {
	rows, err := s.q.RegisterWatcher(ctx, s.watcherPath)
	if err != nil {
		return false, err
	}
	return rows == 1, nil
}

// Load returns the persisted reported set as a filesReported map, with all
// missing-tick counters reset to zero.
func (s *sqliteStore) Load(ctx context.Context) (map[string]int, error) {
	files, err := s.q.ListReportedFiles(ctx, s.watcherPath)
	if err != nil {
		return nil, err
	}
	reported := make(map[string]int, len(files))
	for _, f := range files {
		reported[f] = 0
	}
	return reported, nil
}

func (s *sqliteStore) Add(ctx context.Context, file string) error {
	return s.q.AddReportedFile(ctx, db.AddReportedFileParams{
		WatcherPath: s.watcherPath,
		FilePath:    file,
	})
}

func (s *sqliteStore) Remove(ctx context.Context, file string) error {
	return s.q.RemoveReportedFile(ctx, db.RemoveReportedFileParams{
		WatcherPath: s.watcherPath,
		FilePath:    file,
	})
}
