package testing

import (
	"articles-system/app/configs"
	"articles-system/app/connections"
	"articles-system/app/delivery/routes"
	"articles-system/app/logging"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func SetupTestApp() (logging.LoggerConfig, *sql.DB, *redis.Client, *fiber.App) {
	err := configs.InitConfig()
	if err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}
	appConfig := configs.Config.App

	customLogger := logging.LoggerConfig{LogPath: appConfig.LogPath}
	err = customLogger.InitLogger()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	dbConn, err := connections.InitMySQL()
	if err != nil {
		logging.Error.Fatalf("Error connecting to MySQL: %v", err)
	}

	db, err := dbConn.DB()
	if err != nil {
		logging.Error.Fatalf("Error getting database connection: %v", err)
	}

	redis, err := connections.InitRedis()
	if err != nil {
		logging.Error.Fatalf("Error connecting to Redis: %v", err)
	}

	return customLogger, db, redis, routes.SetUpRouter(dbConn, redis)
}
