package main

import (
	"database/sql"

	mssql "github.com/denisenkom/go-mssqldb"
)

// NewConnector.
func NewConnector(dsn string) (*Connector, error) {
	// Create a new connector object by calling NewConnector
	connector, err := mssql.NewConnector(dsn)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Connector{
		db: db,
	}, nil
}

// Connector.
type Connector struct {
	db      *sql.DB
	Sources Sources
}

// Start.
func (c *Connector) Start() error {
	sources, err := c.getSources()
	if err != nil {
		return err
	}
	c.Sources = sources
	return nil
}

// Close.
func (c *Connector) Close() error {
	return c.db.Close()
}

// Stats.
func (c *Connector) Stats() sql.DBStats {
	return c.db.Stats()
}

func (c *Connector) getSources() (Sources, error) {
	const sql string = `
SELECT DB_Name() AS [databaseName],
SCHEMA_NAME(OBJECTPROPERTY(object_id, 'SchemaId')) AS [schemaName],
OBJECT_NAME(object_id) AS [tableName],
is_track_columns_updated_on AS [trackColumns]
FROM [sys].[change_tracking_tables]`

	rows, err := c.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources Sources
	for rows.Next() {
		source := Source{}
		if err := rows.Scan(&source.Database, &source.Schema, &source.Table, &source.TrackColumns); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return sources, nil
}

// Sources.
type Sources []Source

// Source.
type Source struct {
	Database     string
	Schema       string
	Table        string
	TrackColumns bool
}
