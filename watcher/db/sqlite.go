package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Open opens (creating if needed) the sqlite database at path and runs all
// embedded goose migrations.
func Open(ctx context.Context, path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)"
	sqldb, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db %s: %w", path, err)
	}
	// All watchers share this handle; sqlite is single-writer, and with a single
	// connection there is no SQLITE_BUSY path between them at all.
	sqldb.SetMaxOpenConns(1)

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		sqldb.Close()
		return nil, err
	}
	if err := goose.UpContext(ctx, sqldb, "migrations"); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("migrate sqlite db %s: %w", path, err)
	}
	return sqldb, nil
}
