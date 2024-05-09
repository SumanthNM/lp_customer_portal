/**
* Contains the audit trail for vehicle assignemt
**/

package models

import (
	"time"

	"gorm.io/gorm"
)

type DriverVehicleAssign struct {
	*gorm.Model
	UserID    uint      `json:"driverId" gorm:"default:null"`
	User      User      `gorm:"foriegnKey:UserID" json:"driver"`
	VehicleID uint      `json:"vehicleId"`
	Vehicle   Vehicle   `gorm:"foriegnKey:VehicleID" json:"vehicle"`
	StartDate time.Time `json:"startDate" gorm:"type:timestamp without time zone"`
	EndDate   time.Time `json:"endDate" gorm:"type:timestamp without time zone"`
	LiveLat   float64   `json:"liveLat"`
	LiveLong  float64   `json:"liveLong"`
	PlannerID *uint     `json:"plannerId" gorm:"default:null;foreignKey:PlannerID"`
	IsManual  bool      `json:"isManual"`
}
