package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bcc-code/mediabank-bridge/log"
)

type callbackRequest struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	UpdatedAt time.Time `json:"updatedAt"`
	Size      int64     `json:"size"`
}

// postCallback reports a file to the callback url. It returns an error if the
// notification could not be delivered, so callers can retry on the next tick
// instead of dropping the event.
func postCallback(callbackUrl, path string, file os.FileInfo) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.L.Error().Err(err).Str("file", file.Name()).Msg("Failed to resolve absolute path, using raw path")
		absPath = path
	}

	body, err := json.Marshal(callbackRequest{
		Name:      file.Name(),
		Size:      file.Size(),
		Path:      absPath,
		UpdatedAt: file.ModTime(),
	})
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(callbackUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Drain the body so the keep-alive connection can be reused.
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("callback returned status %d", resp.StatusCode)
	}
	return nil
}
