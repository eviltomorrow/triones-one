package sqlite3

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DSN     string
	MinOpen int = 5
	MaxOpen int = 10
	DB      *sql.DB

	RetryTimes = 3
	Period     = 10 * time.Second
)

var DefaultConnectTimeout = 5 * time.Second

// Build build sqlite3
func Connect() error {
	var (
		pool *sql.DB
		err  error

		i = 1
	)
	for {
		if i > RetryTimes {
			if err != nil {
				return err
			}
			return fmt.Errorf("panic: connect SQLite3 failure, err is nil?")
		}
		pool, err = buildSQLite3(DSN)
		if err == nil {
			break
		}
		if err != nil {
			log.Printf("[E] Try to connect to SQLite3, retry: %d, nest error: %v", i, err)
		}
		i++
		time.Sleep(Period)
	}
	DB = pool

	return nil
}

// Close close sqlite3
func Close() error {
	if DB == nil {
		return nil
	}

	return DB.Close()
}

func buildSQLite3(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("SQLite3: no DSN set")
	}
	pool, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(MaxOpen)
	pool.SetMaxIdleConns(MinOpen)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	if err = pool.PingContext(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

// Exec exec SQLite3
type Exec interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
