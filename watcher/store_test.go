package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bcc-code/bccm-utils/watcher/db"
)

// newTestQueries opens a fresh temp-file database with migrations applied,
// matching how production opens it.
func newTestQueries(t *testing.T) *db.Queries {
	t.Helper()
	sqldb, err := db.Open(context.Background(), filepath.Join(t.TempDir(), "watcher.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqldb.Close() })
	return db.New(sqldb)
}

func TestStoreRegisterAndAdd(t *testing.T) {
	ctx := context.Background()
	q := newTestQueries(t)
	s := newStore(q, "/some/dir/*.mxf")

	firstRun, err := s.Register(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !firstRun {
		t.Fatal("expected first Register to report first run")
	}
	firstRun, err = s.Register(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if firstRun {
		t.Fatal("expected second Register to not report first run")
	}

	// Add is idempotent; Load returns membership with zeroed counters.
	if err := s.Add(ctx, "/some/dir/a.mxf"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add(ctx, "/some/dir/a.mxf"); err != nil {
		t.Fatal(err)
	}
	reported, err := s.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(reported) != 1 {
		t.Fatalf("expected 1 reported file, got %d", len(reported))
	}
	if v, ok := reported["/some/dir/a.mxf"]; !ok || v != 0 {
		t.Fatalf("expected file with counter 0, got %v %v", v, ok)
	}

	// Watchers with a different pattern do not see each other's files.
	other := newStore(q, "/other/dir/*.mxf")
	reported, err = other.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(reported) != 0 {
		t.Fatalf("expected no files for other watcher, got %d", len(reported))
	}

	if err := s.Remove(ctx, "/some/dir/a.mxf"); err != nil {
		t.Fatal(err)
	}
	reported, err = s.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(reported) != 0 {
		t.Fatalf("expected no files after Remove, got %d", len(reported))
	}
}
