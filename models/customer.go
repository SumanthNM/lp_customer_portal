//customer model

package models

import "gorm.io/gorm"

type Customer struct {
	*gorm.Model
	CustomerNo   string `json:"customerId" gorm:"unique"`
	CustomerName string `json:"customerName"`
	Contact      string `json:"contact"`
	Email        string `json:"email"`
}
