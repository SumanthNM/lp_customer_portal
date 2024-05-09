/**
 *
 *
 *
 *
**/package models

import "gorm.io/gorm"

type Zone struct {
	*gorm.Model
	Name           string         `json:"name"`
	ZoneType       string         `json:"zoneType"`
	Status         string         `json:"status"`
	IsActive       bool           `json:"isActive"`
	Description    string         `json:"description"`
	Color          string         `json:"color"`
	ZoneBoundaries []ZoneBoundary `json:"zoneBoundary" validate:"required,dive"`
}

type ZoneBoundary struct {
	*gorm.Model
	ZoneID   int     `json:"zoneId" ` // ZoneID is a Foreign key
	Lat      float64 `json:"lat" validate:"required,float64"`
	Long     float64 `json:"long" validate:"required,float64"`
	Sequence int     `json:"sequence" validate:"required"`
}
