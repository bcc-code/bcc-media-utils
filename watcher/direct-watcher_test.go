package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestDirectWatcher(t *testing.T) {
	var posts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts.Add(1)
	}))
	defer server.Close()

	ctx := context.Background()
	dir := t.TempDir()
	preexisting := filepath.Join(dir, "old.mxf")
	if err := os.WriteFile(preexisting, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "*.mxf")
	store := newStore(newTestQueries(t), path)
	w := &directWatcher{
		path:          path,
		missingTicks:  2,
		filesReported: map[string]int{},
		callbackUrl:   server.URL,
		store:         store,
	}
	// Simulate newWatcher's first-run seeding: preexisting files count as
	// reported and are persisted.
	files, err := filepath.Glob(w.path)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := store.Add(ctx, f); err != nil {
			t.Fatal(err)
		}
		w.filesReported[f] = 0
	}

	w.doWatch(ctx)
	if posts.Load() != 0 {
		t.Fatalf("expected no post for preexisting file, got %d", posts.Load())
	}

	file := filepath.Join(dir, "new.mxf")
	if err := os.WriteFile(file, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	w.doWatch(ctx)
	if posts.Load() != 1 {
		t.Fatalf("expected 1 post for new file, got %d", posts.Load())
	}
	w.doWatch(ctx)
	if posts.Load() != 1 {
		t.Fatalf("expected no duplicate post, got %d", posts.Load())
	}

	// Deleted: forgotten (and the map pruned) only after missingTicks polls.
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	w.doWatch(ctx)
	if _, ok := w.filesReported[file]; !ok {
		t.Fatal("expected file still remembered after 1 missing tick")
	}
	w.doWatch(ctx)
	if _, ok := w.filesReported[file]; ok {
		t.Fatal("expected file pruned after missingTicks")
	}

	// Recreated after being forgotten: reported again.
	if err := os.WriteFile(file, []byte("recreated"), 0644); err != nil {
		t.Fatal(err)
	}
	w.doWatch(ctx)
	if posts.Load() != 2 {
		t.Fatalf("expected 2 posts after recreate, got %d", posts.Load())
	}
}

func TestDirectWatcherCallbackRetry(t *testing.T) {
	var fail atomic.Bool
	var successes atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		successes.Add(1)
	}))
	defer server.Close()

	ctx := context.Background()
	dir := t.TempDir()
	file := filepath.Join(dir, "rec.mxf")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "*.mxf")
	w := &directWatcher{
		path:          path,
		missingTicks:  2,
		filesReported: map[string]int{},
		callbackUrl:   server.URL,
		store:         newStore(newTestQueries(t), path),
	}

	fail.Store(true)
	w.doWatch(ctx)
	w.doWatch(ctx)
	if successes.Load() != 0 {
		t.Fatalf("expected no successful posts while callback fails, got %d", successes.Load())
	}
	if _, reported := w.filesReported[file]; reported {
		t.Fatal("file must not be marked reported while callbacks fail")
	}

	fail.Store(false)
	w.doWatch(ctx)
	if successes.Load() != 1 {
		t.Fatalf("expected 1 successful post after recovery, got %d", successes.Load())
	}

	w.doWatch(ctx)
	if successes.Load() != 1 {
		t.Fatalf("expected no duplicate post, got %d", successes.Load())
	}
}
