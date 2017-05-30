package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/acoshift/session"
)

// New creates new sql store
func New(db *sql.DB, table string) session.Store {
	_, err := db.Exec(fmt.Sprintf(`
		create table if not exists %s (
			k text,
			v blob,
			e timestamp,
			primary key (k),
			index (e)
		);
	`, table))
	if err != nil {
		log.Printf("session: can not create sql table; %v\n", err)
	}
	s := &sqlStore{
		db:              db,
		getQuery:        fmt.Sprintf(`select v, e, now() from %s where k = $1`, table),
		setQuery:        fmt.Sprintf(`insert into %s (k, v, e) values ($1, $2, $3) on conflict (k) do update set v = excluded.v, k = excluded.k`, table),
		delQuery:        fmt.Sprintf(`delete from %s where k = $1`, table),
		expQuery:        fmt.Sprintf(`update %s set e = $2 where k = $1`, table),
		delExpiredQuery: fmt.Sprintf(`delete from %s where e <= now()`, table),
	}
	go s.cleanupWorker()
	return s
}

type sqlStore struct {
	db              *sql.DB
	getQuery        string
	setQuery        string
	delQuery        string
	expQuery        string
	delExpiredQuery string
}

var errNotFound = errors.New("sql: session not found")

func (s *sqlStore) cleanupWorker() {
	time.Sleep(5 * time.Second)
	for {
		s.db.Exec(s.delExpiredQuery)
		time.Sleep(6 * time.Hour)
	}
}

func (s *sqlStore) Get(key string) ([]byte, error) {
	var bs []byte
	var exp *time.Time
	var now time.Time
	err := s.db.QueryRow(s.getQuery, key).Scan(&bs, &exp, &now)
	if err != nil {
		return nil, err
	}
	if exp != nil && exp.Before(now) {
		s.Del(key)
		return nil, errNotFound
	}
	return bs, nil
}

func (s *sqlStore) Set(key string, value []byte, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.db.Exec(s.setQuery, key, value, exp)
	return err
}

func (s *sqlStore) Del(key string) error {
	_, err := s.db.Exec(s.delQuery, key)
	return err
}

func (s *sqlStore) Exp(key string, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.db.Exec(s.expQuery, key, exp)
	return err
}
