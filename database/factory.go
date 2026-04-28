package database

import (
	"fmt"
	"sync"
)

type DatabaseFactory struct {
	mu       sync.RWMutex
	registry map[DbType]DatabaseCreator
	cache    map[string]Database
}

type DatabaseCreator func(config DatabaseConfig) (Database, error)

var (
	factoryInstance *DatabaseFactory
	factoryOnce     sync.Once
)

func GetDatabaseFactory() *DatabaseFactory {
	factoryOnce.Do(func() {
		factoryInstance = &DatabaseFactory{
			registry: make(map[DbType]DatabaseCreator),
			cache:    make(map[string]Database),
		}
		factoryInstance.registerDefaultDrivers()
	})
	return factoryInstance
}

func (f *DatabaseFactory) registerDefaultDrivers() {
	f.registry[MySQL] = NewMySQLDatabase
	f.registry[PostgreSQL] = NewPostgreSQLDatabase
	f.registry[SQLite] = NewSQLiteDatabase
}

func (f *DatabaseFactory) RegisterDriver(dbType DbType, creator DatabaseCreator) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.registry[dbType] = creator
}

func (f *DatabaseFactory) CreateDatabase(config DatabaseConfig) (Database, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	cacheKey := f.getCacheKey(config)

	if cached, ok := f.cache[cacheKey]; ok {
		return cached, nil
	}

	creator, exists := f.registry[config.Type]
	if !exists {
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}

	db, err := creator(config)
	if err != nil {
		return nil, err
	}

	f.cache[cacheKey] = db
	return db, nil
}

func (f *DatabaseFactory) getCacheKey(config DatabaseConfig) string {
	if config.Type == SQLite {
		return fmt.Sprintf("%s:%s", config.Type, config.Path)
	}
	return fmt.Sprintf("%s://%s:%d/%s", config.Type, config.Host, config.Port, config.Database)
}

func (f *DatabaseFactory) CloseAll() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for key, db := range f.cache {
		db.Close()
		delete(f.cache, key)
	}
}
