package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func openPostgreSQL(t *testing.T) *sql.DB {
	t.Helper()

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", user, password, host, port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("can not open postgres database: %v", err)
	}
	return db
}

func TestSQLGetReturnsScanError(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:1/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	s := &SQL{
		DB:           db,
		GetStatement: `select value from sessions where id = $1`,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	data, err := s.Get(ctx, "missing")
	assert.Nil(t, data)
	assert.Error(t, err)
	assert.False(t, errors.Is(err, session.ErrNotFound))
	assert.False(t, errors.Is(err, io.EOF))
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSQL_PostgreSQL(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	db := openPostgreSQL(t)
	defer db.Close()

	db.Exec(`drop table if exists __sql_postgresql`)

	s := (&SQL{DB: db}).
		GeneratePostgreSQLStatement("__sql_postgresql", true).
		GCEvery(50 * time.Millisecond)

	opt := session.StoreOption{TTL: 20 * time.Millisecond}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "a", data, opt)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	b, err := s.Get(ctx, "a")
	assert.Nil(t, b)
	assert.Error(t, err)

	s.Set(ctx, "a", data, opt)
	time.Sleep(100 * time.Millisecond)
	_, err = s.Get(ctx, "a")
	assert.Error(t, err, "expected expired key return error")

	s.Set(ctx, "a", data, session.StoreOption{TTL: time.Second})
	b, err = s.Get(ctx, "a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	s.Del(ctx, "a")
	_, err = s.Get(ctx, "a")
	assert.Error(t, err)
}

func TestSQLWithoutTTL(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	db := openPostgreSQL(t)
	defer db.Close()

	db.Exec(`drop table if exists __sql_postgresql_without_ttl`)

	s := (&SQL{DB: db}).
		GeneratePostgreSQLStatement("__sql_postgresql_without_ttl", true).
		GCEvery(100 * time.Millisecond)

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err := s.Set(ctx, "a", data, opt)
	assert.NoError(t, err)

	b, err := s.Get(ctx, "a")
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
