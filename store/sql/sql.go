package sql

import (
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/session"
)

// New creates new sql store
func New(db *sql.DB, table string) session.Store {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS $1 (
			k TEXT PRIMARY KEY NOT NULL,
			v BLOB,
			e TIMESTAMP
		);
	`, table)
	if err != nil {
		log.Printf("session: can not create sql table; %v\n", err)
	}
	getStmt, _ := db.Prepare(`SELECT v FROM $1 WHERE k = $2;`)
	setStmt, _ := db.Prepare(`
		INSERT INTO $1 (k, v, e)
		VALUES ($2, $3, $4)
		ON CONFLICT (k)
		DO UPDATE SET v = EXCLUDED.v, k = EXCLUDED.k;
	`)
	delStmt, _ := db.Prepare(`DELETE FROM $1 WHERE k = $2;`)
	expStmt, _ := db.Prepare(`UPDATE $1 SET e = $3 WHERE k = $2;`)
	return &sqlStore{db, getStmt, setStmt, delStmt, expStmt, table}
}

type sqlStore struct {
	db      *sql.DB
	getStmt *sql.Stmt
	setStmt *sql.Stmt
	delStmt *sql.Stmt
	expStmt *sql.Stmt
	table   string
}

func (s *sqlStore) Get(key string) ([]byte, error) {
	var bs []byte
	err := s.getStmt.QueryRow(s.table, key).Scan(&bs)
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
	_, err := s.setStmt.Exec(s.table, key, value, exp)
	return err
}

func (s *sqlStore) Del(key string) error {
	_, err := s.delStmt.Exec(s.table, key)
	return err
}

func (s *sqlStore) Exp(key string, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.expStmt.Exec(s.table, key, exp)
	return err
}
