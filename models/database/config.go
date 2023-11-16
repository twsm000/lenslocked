package database

type Config interface {
	DriverName() string
	DataSourceName() string
}
