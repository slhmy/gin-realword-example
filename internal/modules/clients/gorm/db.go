package gorm_client

import (
	"fmt"
	"sync"

	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	initMutx sync.Mutex
	db       *gorm.DB
	driver   string
	name     string
	host     string
	port     string
	username string
	password string
	sslMode  string
)

func init() {
	driver = core.ConfigStore.GetString(shared.ConfigKeyDatabaseDriver)
	if !driverSupported() {
		panic(fmt.Sprintf("unsupported database driver: %s", driver))
	}
	name = core.ConfigStore.GetString(shared.ConfigKeyDatabaseName)
	host = core.ConfigStore.GetString(shared.ConfigKeyDatabaseHost)
	port = core.ConfigStore.GetString(shared.ConfigKeyDatabasePort)
	username = core.ConfigStore.GetString(shared.ConfigKeyDatabaseUsername)
	password = core.ConfigStore.GetString(shared.ConfigKeyDatabasePassword)
	sslMode = core.ConfigStore.GetString(shared.ConfigKeyDatabaseSSLMode)
}

func driverSupported() bool {
	switch driver {
	case "postgres":
		return true
	default:
		return false
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		initMutx.Lock()
		defer initMutx.Unlock()
		if db != nil {
			return db
		}

		var err error
		switch driver {
		case "postgres":
			db, err = openPostgres()
			if err != nil {
				panic(err)
			}
			return db
		default:
			panic(fmt.Sprintf("unsupported database driver: %s", driver))
		}
	}
	return db
}

func openPostgres() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		port,
		username,
		name,
		password,
		sslMode,
	)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}
