package models

import "gorm.io/gorm"

type Role struct {
	*gorm.Model
	Name        string  `json:"name" gorm:"unique"`
	Permissions string  `json:"permissions" `
	Description string  `json:"description"`
	Users       []*User `gorm:"many2many:user_roles;" json:"users"`
}
