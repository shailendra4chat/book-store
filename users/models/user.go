package models

import "gorm.io/gorm"

//User struct declaration
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}
