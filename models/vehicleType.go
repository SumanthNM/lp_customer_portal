/**
*
**/
package models

import "gorm.io/gorm"

type VehicleType struct {
	*gorm.Model
	Name        string `json:"name"`
	Type        string `json:"type"`
	Capacity    int    `json:"capacity"`
	SpeedFactor int    `json:"speedFactor"`
	IsActive    bool   `json:"isActive"`
	VehicleIcon string `json:"vehicle_icon"`
}
