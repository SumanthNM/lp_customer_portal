/**
 *
 *
**/
package models

import (
	"gorm.io/gorm"
)

type Depot struct {
	*gorm.Model
	DepotName        string  `json:"depotName"`
	DepotDescription string  `json:"depotDescription"`
	Latitude         float64 `json:"lat"`
	Longitude        float64 `json:"long"`
	IsActive         bool    `json:"isActive"`
	Address          string  `json:"address"`
}

// BeforeCreate hook to set IsActive to true before creating a new record
func (depot *Depot) BeforeCreate(tx *gorm.DB) (err error) {
	depot.IsActive = true
	return
}
