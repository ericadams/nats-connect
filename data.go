package main

import (
	"database/sql"

	mssql "github.com/denisenkom/go-mssqldb"
)

func NewSource(dsn string) (*Source, error) {
	// Create a new connector object by calling NewConnector
	connector, err := mssql.NewConnector(dsn)
	if err != nil {
		return nil, err
	}
	return &Source{
		db: sql.OpenDB(connector),
	}, nil
}

// Source.
type Source struct {
	db *sql.DB
}

func (s *Source) Close() error {
	return s.db.Close()
}
