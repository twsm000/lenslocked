package postgres

import "fmt"

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
