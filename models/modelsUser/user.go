package modelsuser

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	IsActive  bool   `json:"isActive"`
	Role      string `json:"role"`
}
