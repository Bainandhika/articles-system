package connections

import (
	"articles-system/app/configs"
	"articles-system/lib/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL() (*gorm.DB, error) {
	domainFunc := "[ connections.InitMySQL ]"

	dbConfig := configs.Config.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	fileName := configs.Config.App.LogPath + "/article-system-gorm.log"

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("%s error opening %s: %v", domainFunc, fileName, err)
	}
	defer file.Close()

	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,            // Don't include params in the SQL log
			Colorful:                  false,           // Disable color
		},
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("%s error creating database connection: %v", domainFunc, err)
	}

	// create database if it doesn't exist
	err = dbConn.Exec("CREATE DATABASE IF NOT EXISTS " + dbConfig.Name).Error
    if err!= nil {
        return nil, fmt.Errorf("%s error creating database: %v", domainFunc, err)
    }

	// Automatically migrate the schema to match the struct definitions
	if err = dbConn.AutoMigrate(&models.Article{}); err != nil {
		return nil, fmt.Errorf("%s error migrating database schema: %v", domainFunc, err)
	}

	return dbConn, nil
}
