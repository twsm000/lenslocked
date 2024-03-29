package postgres

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

type Config struct {
	Driver   string `json:"-"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"ssl_mode"`
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

func MigrateFS(db *sql.DB, dir string, fs fs.FS) error {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(fs)
	defer goose.SetBaseFS(nil) // undo the fs change
	return Migrate(db, dir)
}
