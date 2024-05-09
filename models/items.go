// Item Model

package models

import (
	"gorm.io/gorm"
)

type Item struct {
	*gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	SKU         string  `json:"sku" gorm:"unique"`
	Capacity    float64 `json:"capacity"`
	Volume      float64 `json:"volume"`
	Weight      float64 `json:"weight"`
	Status      string  `json:"status"`
	InnerPack   int  `json:"innerPack"`
	// Sort        string  `json:"sort"`
	// Owner       string  `json:"owner"`
	// Type        string  `json:"type"`
}
