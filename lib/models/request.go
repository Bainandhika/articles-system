package models

type AddArticle struct {
	Author string `json:"author" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body string `json:"body" binding:"required"`
}

type GetArticles struct {
	Query string `json:"query"`
	Author string `json:"author"`
}