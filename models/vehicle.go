/**
 *
 *
**/
package models

import "gorm.io/gorm"

type Vehicle struct {
	*gorm.Model
	AssetID       string      `json:"assetId"`
	LicensePlate  string      `json:"licensePlate" gorm:"unique"`
	Name          string      `json:"name"`
	VehicleType   VehicleType `json:"type"`
	VehicleTypeID int         `json:"vehicleTypeId"` // "belongs to" relationship with VehicleType. VehicleTypeID is the foreign key
	DriverName    string      `json:"driverName"`
	Status        string      `json:"status"`
	DepotStartID  uint        `json:"depotStartId"`                              // foreignKey Depot Start
	DepotStart    Depot       `json:"depotStart" gorm:"foreignKey:DepotStartID"` // foreignKey Depot Start
	DepotEndID    uint        `json:"depotEndId"`                                // foreignKey Depot End
	DepotEnd      Depot       `json:"depotEnd" gorm:"foreignKey:DepotEndID"`     // for foreignKey Depot End
	IsActive      bool        `json:"isActive"`
	IsDeleted     bool        `json:"isDeleted"`
	DriverID      uint        `json:"driverId"`
	Driver        User        `gorm:"foreignKey:DriverID" json:"Driver"`
	SkillsetID    uint        `json:"skillsetID"`
	Skillset      Skillset    `json:"skillsets" gorm:"foreignKey:SkillsetID"`
	VfID          uint        `json:"vfId"`
}

// BeforeCreate hook to set IsActive to true before creating a new record
func (vehicle *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	vehicle.IsActive = true
	return
}
