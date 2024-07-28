package services

import (
	"articles-system/app/repositories"
	"articles-system/lib/models"

	"github.com/gofiber/fiber/v2"
)

type articlesService struct {
	repo repositories.ArticlesRepo
}

type ArticlesService interface {
	Create(payload models.AddArticle) *fiber.Error
	GetArticles(queryParams models.GetArticles) ([]models.Article, *fiber.Error)
}

func NewArticlesService(repo repositories.ArticlesRepo) ArticlesService {
    return &articlesService{repo: repo}
}

func (s *articlesService) Create(payload models.AddArticle) *fiber.Error {
	articleData := models.Article{
		Author: payload.Author,
		Title:  payload.Title,
        Body:   payload.Body,
	}

	err := s.repo.Create(articleData)
	if err!= nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

	return nil
}

func (s *articlesService) GetArticles(queryParams models.GetArticles) ([]models.Article, *fiber.Error) {
	articles, err := s.repo.GetArticles(queryParams.Query, queryParams.Author)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return articles, nil
}