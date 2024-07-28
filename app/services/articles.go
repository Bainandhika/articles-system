package services

import (
	"articles-system/app/repositories"
	"articles-system/lib/models"
)

type articlesService struct {
	repo repositories.ArticlesRepo
}

type ArticlesService interface {}

func NewArticlesService(repo repositories.ArticlesRepo) ArticlesService {
    return &articlesService{repo: repo}
}

func (s *articlesService) Create(payload models.AddArticle) error {
	articleData := models.Article{
		Author: payload.Author,
		Title:  payload.Title,
        Body:   payload.Body,
	}

	return s.repo.Create(articleData)
}

func (s *articlesService) GetArticles(queryParams models.GetArticles) ([]models.Article, error) {
	articles, err := s.repo.GetArticles(queryParams.Query, queryParams.Author)
	if err != nil {
		return nil, err
	}

	return articles, nil
}