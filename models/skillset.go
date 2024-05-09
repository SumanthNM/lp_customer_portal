package models

import "gorm.io/gorm"

type Skillset struct {
	*gorm.Model
	IsActive    bool    `json:"isActive"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Remarks     string  `json:"remarks"`
	Category    string  `json:"category"`
	SubCategory string  `json:"subCategory"`
	Users       []*User `gorm:"many2many:user_skillsets;" json:"users"`
}
