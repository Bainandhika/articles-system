package models

type AddArticle struct {
	Author string `json:"author" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body string `json:"body" validate:"required"`
}

type GetArticles struct {
	Query string `json:"query"`
	Author string `json:"author"`
}