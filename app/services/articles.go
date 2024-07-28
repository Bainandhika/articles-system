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
		return ""
	}
}

func (s *articlesService) setArticlesCache(keyParent string, articles []models.Article) {
	domainFunc := "[ services.articlesService.serArticlesCache ]"

	var filteredArticlesID []string
	for _, article := range articles {
		keyID := fmt.Sprintf("article:id:%d", article.ID)
		articleJson, _ := json.Marshal(article)
		err := s.redis.Set(context.Background(), keyID, articleJson, 8*time.Minute).Err()
		if err != nil {
			logging.Error.Printf("%s error setting article cache [key: %s]: %v", domainFunc, keyParent, err)
		}

		if keyParent == "" {
			filteredArticlesID = append(filteredArticlesID, keyID)
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

func (s *articlesService) getArticlesCache(keyParent string) ([]models.Article, error) {
	domainFunc := "[ services.articlesService.getArticlesCache ]"

	if keyParent == "" {
		var (
			cursor uint64
			keys   []string
			err    error
		)

		for {
			var scanKeys []string
			scanKeys, cursor, err = s.redis.Scan(context.Background(), cursor, "*article:id*", 0).Result()
			if err != nil {
				return nil, fmt.Errorf("%s error scanning article cache [key: %s]: %v", domainFunc, keyParent, err)
			}

			keys = append(keys, scanKeys...)

			if cursor == 0 {
				break
			}
		}

		var articles []models.Article
		for i, keyID := range keys {
			articleJson, err := s.redis.Get(context.Background(), keyID).Result()
			if err != nil {
				if i == len(keys)-1 {
					return nil, fmt.Errorf("%s error getting article cache [key: %s]: %v", domainFunc, keyParent, err)
				} else {
					continue
				}
			}

			var article models.Article
			err = json.Unmarshal([]byte(articleJson), &articleJson)
			if err != nil {
				if i == len(keys)-1 {
					return nil, fmt.Errorf("%s error unmarshalling article cache [key: %s]: %v", domainFunc, keyParent, err)
				} else {
					continue
				}
			}

			articles = append(articles, article)
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
	for i, keyID := range articlesKeyID {
		articleJson, err := s.redis.Get(context.Background(), keyID).Result()
		if err != nil {
			if i == len(articlesKeyID)-1 {
				return nil, fmt.Errorf("%s error getting article cache [key: %s]: %v", domainFunc, keyID, err)
			} else {
				continue
			}

		}

		var article models.Article
		err = json.Unmarshal([]byte(articleJson), &article)
		if err != nil {
			if i == len(articlesKeyID)-1 {
				return nil, fmt.Errorf("%s error unmarshalling article cache [key: %s]: %v", domainFunc, keyID, err)
			} else {
				continue
			}
		}

		articles = append(articles, article)
	}

	return articles, nil
}
