package models

import "time"

type Article struct {
	ID      int       `json:"id" gorm:"primary_key;auto_increment"`
	Author  string    `json:"author" gorm:"type:varchar(255);not null"`
	Title   string    `json:"title" gorm:"type:varchar(300);not null"`
	Body    string    `json:"body" gorm:"type:text;not null"`
	Created time.Time `json:"created" gorm:"type:timestamp;not null"`
}
