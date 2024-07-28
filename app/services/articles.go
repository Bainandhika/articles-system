package services

import (
	"articles-system/app/repositories"
	"articles-system/lib/models"
	"fmt"

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
	domainFunc := "[ services.articles.Create ]"

	articleData := models.Article{
		Author: payload.Author,
		Title:  payload.Title,
		Body:   payload.Body,
	}

	err := s.repo.Create(articleData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("%s error creating article: %v", domainFunc, err))
	}

	return nil
}

func (s *articlesService) GetArticles(queryParams models.GetArticles) ([]models.Article, *fiber.Error) {
	domainFunc := "[ services.articles.GetArticles ]"

	articles, err := s.repo.GetArticles(queryParams.Query, queryParams.Author)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("%s error getting article: %v", domainFunc, err))
	}

	return articles, nil
}
