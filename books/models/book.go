package models

import "gorm.io/gorm"

//Book struct declaration
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
	Price  int    `json:"price"`
}
