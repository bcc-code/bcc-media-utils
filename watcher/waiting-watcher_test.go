package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {
	log.ConfigureGlobalLogger(zerolog.DebugLevel)
	os.Exit(m.Run())
}

func TestWaitingWatcher(t *testing.T) {
	var posts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts.Add(1)
	}))
	defer server.Close()

	dir := t.TempDir()
	file := filepath.Join(dir, "rec.mxf")

	w := &waitingWatcher{
		path:          filepath.Join(dir, "*.mxf"),
		stableTicks:   2,
		missingTicks:  2,
		tracked:       map[string]*fileState{},
		filesReported: map[string]int{},
		callbackUrl:   server.URL,
	}

	write := func(content string) {
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// setModTime makes each write's mtime distinct even on coarse filesystems.
	setModTime := func(tm time.Time) {
		if err := os.Chtimes(file, tm, tm); err != nil {
			t.Fatal(err)
		}
	}

	now := time.Now()

	write("a")
	setModTime(now)
	w.doWatch()
	if posts.Load() != 0 {
		t.Fatalf("expected no post after first sight, got %d", posts.Load())
	}

	write("ab")
	setModTime(now.Add(time.Second))
	w.doWatch()
	if posts.Load() != 0 {
		t.Fatalf("expected no post while growing, got %d", posts.Load())
	}

	// Unchanged for one tick — below stableTicks threshold.
	w.doWatch()
	if posts.Load() != 0 {
		t.Fatalf("expected no post after 1 stable tick, got %d", posts.Load())
	}

	// Second stable tick — should report exactly once.
	w.doWatch()
	if posts.Load() != 1 {
		t.Fatalf("expected 1 post after %d stable ticks, got %d", w.stableTicks, posts.Load())
	}

	// Further ticks must not repost.
	w.doWatch()
	w.doWatch()
	if posts.Load() != 1 {
		t.Fatalf("expected no duplicate posts, got %d", posts.Load())
	}

	// Transient disappearance (shorter than missingTicks) must not repost.
	hidden := file + ".hidden"
	if err := os.Rename(file, hidden); err != nil {
		t.Fatal(err)
	}
	w.doWatch()
	if err := os.Rename(hidden, file); err != nil {
		t.Fatal(err)
	}
	for range 5 {
		w.doWatch()
	}
	if posts.Load() != 1 {
		t.Fatalf("expected no repost after transient disappearance, got %d", posts.Load())
	}

	// Genuine delete: gone for >= missingTicks, then recreated and stabilized -> second post.
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	w.doWatch()
	w.doWatch()
	if _, reported := w.filesReported[file]; reported {
		t.Fatal("expected file to be forgotten after missingTicks")
	}

	write("new recording")
	setModTime(now.Add(10 * time.Second))
	w.doWatch() // first sight
	w.doWatch() // stable tick 1
	w.doWatch() // stable tick 2 -> report
	if posts.Load() != 2 {
		t.Fatalf("expected 2 posts after recreate and stabilize, got %d", posts.Load())
	}
}

func TestWaitingWatcherCallbackRetry(t *testing.T) {
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

	dir := t.TempDir()
	file := filepath.Join(dir, "rec.mxf")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	w := &waitingWatcher{
		path:          filepath.Join(dir, "*.mxf"),
		stableTicks:   1,
		missingTicks:  2,
		tracked:       map[string]*fileState{},
		filesReported: map[string]int{},
		callbackUrl:   server.URL,
	}

	fail.Store(true)
	w.doWatch() // first sight
	w.doWatch() // stable -> report attempt fails
	w.doWatch() // retried, still failing
	if successes.Load() != 0 {
		t.Fatalf("expected no successful posts while callback fails, got %d", successes.Load())
	}
	if _, reported := w.filesReported[file]; reported {
		t.Fatal("file must not be marked reported while callbacks fail")
	}

	// Endpoint recovers: the next tick delivers the event exactly once.
	fail.Store(false)
	w.doWatch()
	if successes.Load() != 1 {
		t.Fatalf("expected 1 successful post after recovery, got %d", successes.Load())
	}
	if _, reported := w.filesReported[file]; !reported {
		t.Fatal("file must be marked reported after successful callback")
	}

	w.doWatch()
	if successes.Load() != 1 {
		t.Fatalf("expected no duplicate post, got %d", successes.Load())
	}
}
