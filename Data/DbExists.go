package data

import (
	"database/sql"
	"log"
	"os"
)

func DbExists(dbPath string) {
	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if !dbExists {
		InitDB(db)
	}
}
