package testdata

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "test_habitmap.db"

type TestDB struct {
	DB *sql.DB
}

func NewTestDB() TestDB {
	file, err := os.Create(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	return TestDB{DB: db}
}

func (t *TestDB) Setup() {
	script, err := os.ReadFile("../sqlite/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = t.DB.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	// TODO: seed database
}

func (t *TestDB) Teardown() {
	t.DB.Close()
}
