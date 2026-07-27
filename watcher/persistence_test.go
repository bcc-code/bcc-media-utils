package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// TestWaitingWatcherRestartDedup verifies that a file reported before a
// restart is not reported again by a fresh watcher loading the same database.
func TestWaitingWatcherRestartDedup(t *testing.T) {
	var posts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts.Add(1)
	}))
	defer server.Close()

	ctx := context.Background()
	dir := t.TempDir()
	file := filepath.Join(dir, "rec.mxf")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "*.mxf")
	q := newTestQueries(t)

	first := &waitingWatcher{
		path:          path,
		stableTicks:   1,
		missingTicks:  2,
		tracked:       map[string]*fileState{},
		filesReported: map[string]int{},
		callbackUrl:   server.URL,
		store:         newStore(q, path),
	}
	first.doWatch(ctx) // first sight
	first.doWatch(ctx) // stable -> report
	if posts.Load() != 1 {
		t.Fatalf("expected 1 post before restart, got %d", posts.Load())
	}

	// "Restart": a fresh watcher on the same database, seeded like newWatcher.
	store := newStore(q, path)
	reported, err := store.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	second := &waitingWatcher{
		path:          path,
		stableTicks:   1,
		missingTicks:  2,
		tracked:       map[string]*fileState{},
		filesReported: reported,
		callbackUrl:   server.URL,
		store:         store,
	}
	for range 5 {
		second.doWatch(ctx)
	}
	if posts.Load() != 1 {
		t.Fatalf("expected no repost after restart, got %d", posts.Load())
	}

	// A new file after the restart is still reported.
	newFile := filepath.Join(dir, "new.mxf")
	if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := os.Chtimes(newFile, now, now); err != nil {
		t.Fatal(err)
	}
	second.doWatch(ctx) // first sight
	second.doWatch(ctx) // stable -> report
	if posts.Load() != 2 {
		t.Fatalf("expected 2 posts after new file, got %d", posts.Load())
	}
}

// TestDirectWatcherBackfill verifies the first-run seed (no callbacks for
// preexisting files) and that files appearing while the watcher is down are
// reported after a restart.
func TestDirectWatcherBackfill(t *testing.T) {
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
	q := newTestQueries(t)

	// First run: register + seed, like newWatcher does.
	store := newStore(q, path)
	firstRun, err := store.Register(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !firstRun {
		t.Fatal("expected first run")
	}
	reported := map[string]int{}
	files, err := filepath.Glob(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := store.Add(ctx, f); err != nil {
			t.Fatal(err)
		}
		reported[f] = 0
	}
	first := &directWatcher{
		path:          path,
		missingTicks:  2,
		filesReported: reported,
		callbackUrl:   server.URL,
		store:         store,
	}
	first.doWatch(ctx)
	if posts.Load() != 0 {
		t.Fatalf("expected no post for seeded preexisting file, got %d", posts.Load())
	}

	// Watcher is "down": a new file appears.
	missed := filepath.Join(dir, "missed.mxf")
	if err := os.WriteFile(missed, []byte("missed"), 0644); err != nil {
		t.Fatal(err)
	}

	// Restart: not a first run, no snapshot — the loaded set decides.
	store2 := newStore(q, path)
	firstRun, err = store2.Register(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if firstRun {
		t.Fatal("expected not a first run after restart")
	}
	loaded, err := store2.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded[preexisting]; !ok {
		t.Fatal("expected preexisting file in loaded set")
	}
	second := &directWatcher{
		path:          path,
		missingTicks:  2,
		filesReported: loaded,
		callbackUrl:   server.URL,
		store:         store2,
	}
	second.doWatch(ctx)
	if posts.Load() != 1 {
		t.Fatalf("expected exactly 1 post backfilling the missed file, got %d", posts.Load())
	}
	second.doWatch(ctx)
	if posts.Load() != 1 {
		t.Fatalf("expected no duplicate post, got %d", posts.Load())
	}
}

// TestDirectWatcherRemovalPersisted verifies that aging a file out of the
// reported set also removes it from the database, so a recreated file is
// reported again even across a restart.
func TestDirectWatcherRemovalPersisted(t *testing.T) {
	var posts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts.Add(1)
	}))
	defer server.Close()

	ctx := context.Background()
	dir := t.TempDir()
	file := filepath.Join(dir, "rec.mxf")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
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

	w.doWatch(ctx) // report
	if posts.Load() != 1 {
		t.Fatalf("expected 1 post, got %d", posts.Load())
	}
	reported, err := store.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reported[file]; !ok {
		t.Fatal("expected reported file persisted")
	}

	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	w.doWatch(ctx) // missing tick 1
	w.doWatch(ctx) // missing tick 2 -> forgotten
	reported, err = store.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reported[file]; ok {
		t.Fatal("expected aged-out file removed from store")
	}
}
