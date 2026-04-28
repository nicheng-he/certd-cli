package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	*BaseDatabase
	config DatabaseConfig
}

func NewMySQLDatabase(config DatabaseConfig) (Database, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	if len(config.Options) > 0 {
		dsn += "&" + strings.Join(config.Options, "&")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开 MySQL 连接失败: %w", err)
	}

	baseDb := NewBaseDatabase(db)
	baseDb.configureConnection(db)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("连接 MySQL 数据库失败: %w", err)
	}

	log.Println("MySQL 数据库连接成功")
	return &MySQLDatabase{
		BaseDatabase: baseDb,
		config:       config,
	}, nil
}
