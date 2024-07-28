package routes

import (
	"articles-system/app/handlers"
	"articles-system/app/repositories"
	"articles-system/app/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUpRouter(db *gorm.DB) {
	repo := repositories.NewArticlesRepo(db)
	service := services.NewArticlesService(repo)
	handlers := handlers.NewArticlesHandler(service)

	r := fiber.New()
	api := r.Group("/articles")
	api.Post("/", handlers.Create)
	api.Get("/", handlers.GetArticles)
}