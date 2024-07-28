package models

type AddArticle struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body string `json:"body"`
}

type GetArticles struct {
	Query string `json:"query"`
	Author string `json:"author"`
}