package app

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitSQLite(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS queues (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        song_name TEXT
    );
    `
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return DB.Ping()
}

func Close() error {
	return DB.Close()
}
