package main

import (
	"articles-system/app/configs"
	"articles-system/app/connections"
	"articles-system/app/delivery/routes"
	"articles-system/app/logging"
	"fmt"
	"log"
)

func main() {
	err := configs.InitConfig()
	if err!= nil {
        log.Fatalf("Error initializing configuration: %v", err)
    }
	appConfig := configs.Config.App

	customLogger := logging.LoggerConfig{LogPath: appConfig.LogPath}
	err = customLogger.InitLogger()
	if err!= nil {
        log.Fatalf("Error initializing logger: %v", err)
    }
	defer customLogger.Close()

	dbConn, err := connections.InitMySQL()
	if err!= nil {
        logging.Error.Fatalf("Error connecting to MySQL: %v", err)
    }

	db, err := dbConn.DB()
	if err != nil {
		logging.Error.Fatalf("Error getting database connection: %v", err)
	}
	defer db.Close()

	redis, err := connections.InitRedis()
	if err!= nil {
        logging.Error.Fatalf("Error connecting to Redis: %v", err)
    }
	defer redis.Close()

	router := routes.SetUpRouter(dbConn, redis)

	address := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	if err := router.Listen(address); err != nil {
		logging.Error.Fatalf("Error listening to %s: %v", address, err)
	}
}