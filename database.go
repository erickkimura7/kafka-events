package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConfigDB() {
	connStr := "postgresql://postgres:postgres@localhost:5432/todos?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(db)
}
