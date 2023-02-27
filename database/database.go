package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"faceit/appdata"
	"faceit/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	User   string
	Pass   string
	Addr   string
	Port   string
	DbName string
}

var (
	conn *gorm.DB
	m    sync.Mutex
)

func GetConnection() (*gorm.DB, error) {
	if conn == nil {
		m.Lock()
		defer m.Unlock()

		if conn != nil {
			return conn, nil
		}

		config := Config{
			User:   appdata.Config.Database.User,
			Pass:   appdata.Config.Database.Password,
			Addr:   appdata.Config.Database.Host,
			Port:   appdata.Config.Database.Port,
			DbName: appdata.Config.Database.DbName,
		}

		var err error
		conn, err = newDbConnection(config)
		if err != nil {
			return nil, err
		}

	}

	return conn, nil
}

func EnsureDbSchema() error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}

func newDbConnection(dbfig Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbfig.User,
		dbfig.Pass,
		dbfig.Addr,
		dbfig.Port,
		dbfig.DbName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                  // Slow SQL threshold
			LogLevel:                  logger.Silent,                // Log level
			IgnoreRecordNotFoundError: appdata.Config.Env == "prod", // Ignore ErrRecordNotFound error for logger
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
