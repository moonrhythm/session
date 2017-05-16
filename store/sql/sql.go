package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/acoshift/session"
)

// New creates new sql store
func New(db *sql.DB, table string) session.Store {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			k TEXT PRIMARY KEY NOT NULL,
			v BLOB,
			e TIMESTAMP,
			INDEX (e)
		);
	`, table))
	if err != nil {
		log.Printf("session: can not create sql table; %v\n", err)
	}
	getStmt, _ := db.Prepare(fmt.Sprintf(`SELECT v FROM %s WHERE k = $1;`, table))
	setStmt, _ := db.Prepare(fmt.Sprintf(`
		INSERT INTO %s (k, v, e)
		VALUES ($1, $2, $3)
		ON CONFLICT (k)
		DO UPDATE SET v = EXCLUDED.v, k = EXCLUDED.k;
	`, table))
	delStmt, _ := db.Prepare(fmt.Sprintf(`DELETE FROM %s WHERE k = $1;`, table))
	expStmt, _ := db.Prepare(fmt.Sprintf(`UPDATE %s SET e = $2 WHERE k = $1;`, table))
	return &sqlStore{db, getStmt, setStmt, delStmt, expStmt}
}

type sqlStore struct {
	db      *sql.DB
	getStmt *sql.Stmt
	setStmt *sql.Stmt
	delStmt *sql.Stmt
	expStmt *sql.Stmt
}

func (s *sqlStore) Get(key string) ([]byte, error) {
	var bs []byte
	err := s.getStmt.QueryRow(key).Scan(&bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (s *sqlStore) Set(key string, value []byte, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.setStmt.Exec(key, value, exp)
	return err
}

func (s *sqlStore) Del(key string) error {
	_, err := s.delStmt.Exec(key)
	return err
}

func (s *sqlStore) Exp(key string, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.expStmt.Exec(key, exp)
	return err
}
