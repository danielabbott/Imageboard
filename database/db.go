package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type ImageID int32
type TagID int32

type DB struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func get_statement(db *DB, sql string) *sql.Stmt {
	if existing, ok := db.statements[sql]; ok {
		return existing
	} else {
		st, err := db.db.Prepare(sql)

		if err != nil {
			log.Fatal(err)
		}

		return st
	}
}

func InitDatabase() DB {
	var connStr = "user=postgres dbname=imageboard password='p123' sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return DB{
		db: db,
	}
}
