package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Config struct {
	Driver   string
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c Config) DriverName() string {
	return c.Driver
}

func (c Config) DataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSLMode,
	)
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, dir); err != nil {
		return err
	}

	return nil
}
