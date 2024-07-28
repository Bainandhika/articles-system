package services

import (
	"articles-system/app/logging"
	"articles-system/app/repositories"
	"articles-system/lib/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type articlesService struct {
	redis *redis.Client
	repo  repositories.ArticlesRepo
}

type ArticlesService interface {
	Create(payload models.AddArticle) *fiber.Error
	GetArticles(queryParams models.GetArticles) ([]models.Article, *fiber.Error)
}

func NewArticlesService(redis *redis.Client, repo repositories.ArticlesRepo) ArticlesService {
	return &articlesService{
		redis: redis,
		repo:  repo,
	}
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

	keyParent := s.createKeyParent(queryParams.Query, queryParams.Author)
	articles, err := s.getArticlesCache(keyParent)
	if err != nil {
		logging.Error.Printf("%s error getting article: %v", domainFunc, err)

		articles, err = s.repo.GetArticles(queryParams.Query, queryParams.Author)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("%s error getting article: %v", domainFunc, err))
		}
	}

	go s.setArticlesCache(keyParent, articles)

	return articles, nil
}

func (s *articlesService) createKeyParent(query, author string) string {
	if query != "" && author != "" {
		return fmt.Sprintf("articles:query:%s:author:%s", query, author)
	} else if query != "" {
		return fmt.Sprintf("articles:query:%s", query)
	} else if author != "" {
		return fmt.Sprintf("articles:author:%s", author)
	} else {
		return "articles:all"
	}
}

func (s *articlesService) getArticlesCache(keyParent string) ([]models.Article, error) {
	domainFunc := "[ services.articlesService.getArticlesCache ]"

	if keyParent == "articles:all" {
		articlesJson, err := s.redis.Get(context.Background(), keyParent).Result()
		if err != nil {
			return nil, fmt.Errorf("%s error getting article cache [key: %s]: %v", domainFunc, keyParent, err)
		}

		var articles []models.Article
		err = json.Unmarshal([]byte(articlesJson), &articles)
		if err != nil {
			return nil, fmt.Errorf("%s error unmarshalling article cache [key: %s]: %v", domainFunc, keyParent, err)
		}

		return articles, nil
	}

	articlesKeyIDJson, err := s.redis.Get(context.Background(), keyParent).Result()
	if err != nil {
		return nil, fmt.Errorf("%s error getting article cache [key: %s]: %v", domainFunc, keyParent, err)
	}

	var articlesKeyID []string
	err = json.Unmarshal([]byte(articlesKeyIDJson), &articlesKeyID)
	if err != nil {
		return nil, fmt.Errorf("%s error unmarshalling article cache [key: %s]: %v", domainFunc, keyParent, err)
	}

	var articles []models.Article
	for _, keyID := range articlesKeyID {
		articleJson, err := s.redis.Get(context.Background(), keyID).Result()
		if err != nil {
			return nil, fmt.Errorf("%s error getting article cache [key: %s]: %v", domainFunc, keyID, err)
		}

		var article models.Article
		err = json.Unmarshal([]byte(articleJson), &article)
		if err != nil {
			return nil, fmt.Errorf("%s error unmarshalling article cache [key: %s]: %v", domainFunc, keyID, err)
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *articlesService) setArticlesCache(keyParent string, articles []models.Article) {
	domainFunc := "[ services.articlesService.serArticlesCache ]"

	if keyParent == "articles:all" {
		articlesJson, _ := json.Marshal(articles)
		err := s.redis.Set(context.Background(), keyParent, articlesJson, 6*time.Minute).Err()
		if err != nil {
			logging.Error.Printf("%s error setting article cache [key: %s]: %v", domainFunc, keyParent, err)
		}
	} else {
		var filteredArticlesID []string
		for _, article := range articles {
			keyID := fmt.Sprintf("article:%d", article.ID)
			filteredArticlesID = append(filteredArticlesID, keyID)

			articleJson, _ := json.Marshal(article)
			err := s.redis.Set(context.Background(), keyID, articleJson, 6*time.Minute).Err()
			if err != nil {
				logging.Error.Printf("%s error setting article cache [key: %s]: %v", domainFunc, keyID, err)
			}
		}

		if len(filteredArticlesID) > 0 {
			filteredArticlesIDJson, _ := json.Marshal(filteredArticlesID)
			err := s.redis.Set(context.Background(), keyParent, filteredArticlesIDJson, 5*time.Minute).Err()
			if err != nil {
				logging.Error.Printf("%s error setting article cache [key: %s]: %v", domainFunc, keyParent, err)
			}
		}
	}
}
