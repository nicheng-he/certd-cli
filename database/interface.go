package database

import "database/sql"

type DbType string

const (
	MySQL      DbType = "mysql"
	PostgreSQL DbType = "postgres"
	SQLite     DbType = "sqlite"
)

func (dt DbType) String() string {
	return string(dt)
}

type DatabaseConfig struct {
	Type     DbType   `mapstructure:"type" yaml:"type" json:"type"`
	Host     string   `mapstructure:"host" yaml:"host" json:"host"`
	Port     int      `mapstructure:"port" yaml:"port" json:"port"`
	Username string   `mapstructure:"username" yaml:"username" json:"username"`
	Password string   `mapstructure:"password" yaml:"password" json:"password"`
	Database string   `mapstructure:"database" yaml:"database" json:"database"`
	Path     string   `mapstructure:"path" yaml:"path" json:"path"`
	Options  []string `mapstructure:"options" yaml:"options" json:"options"`
}

type Database interface {
	GetDB() *sql.DB
	Ping() error
	Close() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
}
