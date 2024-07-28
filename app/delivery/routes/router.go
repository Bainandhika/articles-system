package routes

import (
	"articles-system/app/delivery/handlers"
	"articles-system/app/repositories"
	"articles-system/app/services"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUpRouter(db *gorm.DB, redis *redis.Client) *fiber.App {
	repo := repositories.NewArticlesRepo(db)
	service := services.NewArticlesService(redis, repo)
	handlers := handlers.NewArticlesHandler(service)

	router := fiber.New()
	api := router.Group("/articles")
	api.Post("/", handlers.Create)
	api.Get("/", handlers.GetArticles)

	return router
}