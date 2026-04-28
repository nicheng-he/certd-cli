package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type PostgreSQLDatabase struct {
	*BaseDatabase
	config DatabaseConfig
}

func NewPostgreSQLDatabase(config DatabaseConfig) (Database, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)

	if len(config.Options) > 0 {
		dsn += " " + strings.Join(config.Options, " ")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开 PostgreSQL 连接失败: %w", err)
	}

	baseDb := NewBaseDatabase(db)
	baseDb.configureConnection(db)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("连接 PostgreSQL 数据库失败: %w", err)
	}

	log.Println("PostgreSQL 数据库连接成功")
	return &PostgreSQLDatabase{
		BaseDatabase: baseDb,
		config:       config,
	}, nil
}
