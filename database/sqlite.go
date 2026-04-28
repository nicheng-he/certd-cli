package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type SQLiteDatabase struct {
	*BaseDatabase
	config DatabaseConfig
}

func NewSQLiteDatabase(config DatabaseConfig) (Database, error) {
	// 前置处理：如果 Path 为空但 Database 有值，将 Database 的值赋给 Path
	if config.Path == "" && config.Database != "" {
		config.Path = config.Database
		config.Database = ""
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	path := config.Path
	if path == "" {
		path = "data.db"
	}

	if _, err := os.Stat(path); err == nil {
		info, _ := os.Stat(path)
		fmt.Printf("%s文件存在，大小: %d 字节\n", path, info.Size())
	} else {
		fmt.Printf("%s文件不存在: %v\n", path, err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("打开 SQLite 连接失败: %w", err)
	}

	baseDb := NewBaseDatabase(db)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("连接 SQLite 数据库失败: %w", err)
	}

	log.Println("SQLite 数据库连接成功")
	return &SQLiteDatabase{
		BaseDatabase: baseDb,
		config:       config,
	}, nil
}
