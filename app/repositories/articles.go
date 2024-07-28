package repositories

import (
	"articles-system/lib/models"

	"gorm.io/gorm"
)

type articlesRepo struct {
	db *gorm.DB
}

type ArticlesRepo interface {
	Create(data models.Article) error
	GetArticles(query, author string) ([]models.Article, error)
}

func NewArticlesRepo(db *gorm.DB) ArticlesRepo {
    return &articlesRepo{db: db}
}

func (r *articlesRepo) Create(data models.Article) error {
	return r.db.Create(&data).Error
}

func (r *articlesRepo) GetArticles(query, author string) ([]models.Article, error) {
	var articles []models.Article

    if query != "" || author != "" {
        query = "%" + query + "%"
        r.db.Where("title LIKE ? OR body LIKE ? OR author = ?", query, query, author).Find(&articles)
    } else {
        r.db.Find(&articles)
    }

    return articles, nil
}