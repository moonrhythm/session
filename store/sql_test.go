package store

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/moonrhythm/session"
)

func TestSQL_PostgreSQL(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("postgres", "postgres://localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("can not open postgres database: %v", err)
	}
	defer db.Close()

	db.Exec(`drop table if exists __sql_postgresql`)

	s := (&SQL{DB: db}).
		GeneratePostgreSQLStatement("__sql_postgresql", true).
		GCEvery(50 * time.Millisecond)

	opt := session.StoreOption{TTL: 20 * time.Millisecond}

	data := make(session.Data)
	data["test"] = "123"

	err = s.Set("a", data, opt)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	b, err := s.Get("a", opt)
	assert.Nil(t, b)
	assert.Error(t, err)

	s.Set("a", data, opt)
	time.Sleep(100 * time.Millisecond)
	_, err = s.Get("a", opt)
	assert.Error(t, err, "expected expired key return error")

	s.Set("a", data, session.StoreOption{TTL: time.Second})
	b, err = s.Get("a", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)

	_, err = s.Get("a", session.StoreOption{Rolling: true, TTL: time.Minute})
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = s.Get("a", opt)
	assert.NoError(t, err)

	s.Del("a", opt)
	_, err = s.Get("a", opt)
	assert.Error(t, err)
}

func TestSQLWithoutTTL(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("postgres", "postgres://localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("can not open postgres database: %v", err)
	}
	defer db.Close()

	db.Exec(`drop table if exists __sql_postgresql_without_ttl`)

	s := (&SQL{DB: db}).
		GeneratePostgreSQLStatement("__sql_postgresql_without_ttl", true).
		GCEvery(100 * time.Millisecond)

	opt := session.StoreOption{}

	data := make(session.Data)
	data["test"] = "123"

	err = s.Set("a", data, opt)
	assert.NoError(t, err)

	b, err := s.Get("a", opt)
	assert.NoError(t, err)
	assert.Equal(t, data, b)
}
