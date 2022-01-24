package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Datastore for Postgres
type Datastore struct {
	DB *sql.DB
}

// NewDataStore creates new Postgres data store
func NewDataStore(dsn string) *Datastore {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	return &Datastore{
		DB: db,
	}
}
