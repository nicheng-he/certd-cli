package database

import (
	"database/sql"
	"fmt"
)

type BaseDatabase struct {
	db *sql.DB
}

func NewBaseDatabase(db *sql.DB) *BaseDatabase {
	return &BaseDatabase{db: db}
}

func (b *BaseDatabase) GetDB() *sql.DB {
	return b.db
}

func (b *BaseDatabase) Ping() error {
	return b.db.Ping()
}

func (b *BaseDatabase) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

func (b *BaseDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return b.db.Query(query, args...)
}

func (b *BaseDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return b.db.QueryRow(query, args...)
}

func (b *BaseDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return b.db.Exec(query, args...)
}

func (b *BaseDatabase) Begin() (*sql.Tx, error) {
	return b.db.Begin()
}

func (b *BaseDatabase) configureConnection(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(300)
}

func validateConfig(config DatabaseConfig) error {
	if config.Type == SQLite {
		if config.Path == "" {
			return fmt.Errorf("SQLite 数据库路径不能为空")
		}
	} else {
		if config.Host == "" {
			return fmt.Errorf("数据库主机不能为空")
		}
		if config.Port == 0 {
			return fmt.Errorf("数据库端口不能为0")
		}
		if config.Database == "" {
			return fmt.Errorf("数据库名称不能为空")
		}
	}
	return nil
}
