package database

import (
	"fmt"
	"sync"
)

type Container struct {
	mu          sync.RWMutex
	services    map[string]interface{}
	databases   map[string]Database
	config      *AppConfig
	initialized bool
}

type AppConfig struct {
	Database DatabaseConfig `mapstructure:"database" yaml:"database" json:"database"`
}

var (
	containerInstance *Container
	containerOnce     sync.Once
)

func GetContainer() *Container {
	containerOnce.Do(func() {
		containerInstance = &Container{
			services:  make(map[string]interface{}),
			databases: make(map[string]Database),
		}
	})
	return containerInstance
}

func (c *Container) Initialize(config *AppConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("容器已初始化")
	}

	c.config = config

	if config.Database.Type != "" {
		db, err := GetDatabaseFactory().CreateDatabase(config.Database)
		if err != nil {
			return fmt.Errorf("初始化数据库失败: %w", err)
		}
		c.databases["default"] = db
		c.services["database"] = db
	}

	c.initialized = true
	return nil
}

func (c *Container) GetDatabase(name ...string) (Database, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dbName := "default"
	if len(name) > 0 {
		dbName = name[0]
	}

	db, exists := c.databases[dbName]
	if !exists {
		return nil, fmt.Errorf("数据库 '%s' 未找到", dbName)
	}

	return db, nil
}

func (c *Container) GetService(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("服务 '%s' 未找到", name)
	}

	return service, nil
}

func (c *Container) RegisterService(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

func (c *Container) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

func (c *Container) Shutdown() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, db := range c.databases {
		db.Close()
		delete(c.databases, name)
	}

	c.services = make(map[string]interface{})
	c.initialized = false
}
