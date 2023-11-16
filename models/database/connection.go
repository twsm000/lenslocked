package database

import "database/sql"

// NewConnection returns an *sql.DB connection open and validated (calls Ping)
func NewConnection(config Config) (*sql.DB, error) {
	db, err := sql.Open(config.DriverName(), config.DataSourceName())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
