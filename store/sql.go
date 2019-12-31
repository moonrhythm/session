package store

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/moonrhythm/session"
)

// SQL is the sql store
type SQL struct {
	DB *sql.DB

	SetStatement string
	GetStatement string
	DelStatement string
	GCStatement  string
}

const (
	pgsqlInitSchema = `create table if not exists %s (
    id varchar,
    value bytea not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz,
    primary key (id)
);
create index if not exists %s_expires_at_idx on %s (expires_at);`
	pgsqlSet = `insert into %s (id, value, created_at, expires_at)
values ($1, $2, $3, $4)
on conflict (id) do update
set value = excluded.value,
    expires_at = excluded.expires_at`
	pgsqlGet = `select value from %s where id = $1 and (expires_at is null or expires_at > now())`
	pgsqlDel = `delete from %s where id = $1`
	pgsqlGC  = `delete from %s where expires_at <= now()`
)

// GeneratePostgrSQLStatement generates postgresql statement
func (s *SQL) GeneratePostgreSQLStatement(table string, initSchema bool) *SQL {
	if initSchema {
		q := fmt.Sprintf(pgsqlInitSchema, table, table, table)
		_, err := s.DB.Exec(q)
		if err != nil {
			log.Printf("store/sql: init postgresql schema error: %v", err)
		}
	}

	s.SetStatement = fmt.Sprintf(pgsqlSet, table)
	s.GetStatement = fmt.Sprintf(pgsqlGet, table)
	s.DelStatement = fmt.Sprintf(pgsqlDel, table)
	s.GCStatement = fmt.Sprintf(pgsqlGC, table)
	return s
}

// Get gets session data from sql db
func (s *SQL) Get(key string, opt session.StoreOption) (session.Data, error) {
	var b []byte
	err := s.DB.QueryRow(s.GetStatement, key).Scan(&b)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, session.ErrNotFound
	}

	d := s.decode(b)

	if opt.Rolling {
		s.Set(key, d, opt)
	}

	return d, nil
}

// Set sets session data to sql db
func (s *SQL) Set(key string, value session.Data, opt session.StoreOption) error {
	now := time.Now()
	var exp sql.NullTime
	if opt.TTL > 0 {
		exp.Valid = true
		exp.Time = now.Add(opt.TTL)
	}

	_, err := s.DB.Exec(s.SetStatement, key, s.encode(value), now, exp)
	return err
}

// Del deletes session data from sql db
func (s *SQL) Del(key string, opt session.StoreOption) error {
	_, err := s.DB.Exec(s.DelStatement, key)
	return err
}

// GC runs gc
func (s *SQL) GC() error {
	_, err := s.DB.Exec(s.GCStatement)
	return err
}

func (s *SQL) gcWorker(d time.Duration) {
	s.GC()
	time.AfterFunc(d, func() { s.gcWorker(d) })
}

// GCEvery runs gc every given duration
func (s *SQL) GCEvery(d time.Duration) *SQL {
	time.AfterFunc(d, func() { s.gcWorker(d) })
	return s
}

func (s *SQL) encode(v interface{}) []byte {
	var buf bytes.Buffer
	session.DefaultStoreCoder.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}

func (s *SQL) decode(b []byte) session.Data {
	var d session.Data
	session.DefaultStoreCoder.NewDecoder(bytes.NewReader(b)).Decode(&d)
	return d
}
