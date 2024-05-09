// address

package models

import "gorm.io/gorm"

type Address struct {
	*gorm.Model
	BuildingNo  string `json:"buildingNo"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ZipCode     uint `json:"zipCode"`
	Address_Str string `json:"addressStr" gorm:"unique"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ZoneID    uint    `json:"zoneId" gorm:"default:null"`
	Zone      *Zone   `gorm:"foreignKey:ZoneID" json:"zone"`
}
