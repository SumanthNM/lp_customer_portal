package models

import (
	"time"

	"gorm.io/gorm"
)

type Planner struct {
	*gorm.Model
	PlannedDate                time.Time             `json:"plannedDate"`
	PlannedBy                  *uint                 `json:"plannedBy" gorm:"default:null"`
	Planned                    User                  `gorm:"foreignKey:PlannedBy" json:"planned"`
	Objective                  string                `json:"objective"`
	Constraints                string                `json:"constraints"`
	Rules                      string                `json:"rules"`
	Date                       string                `json:"date"`
	Jobs                       []Jobs                `json:"jobs"`
	StartedAt                  time.Time             `json:"startedAt"`
	EndedAt                    time.Time             `json:"endedAt"`
	IsSuccess                  bool                  `json:"isSuccess"`
	Lock                       bool                  `json:"lock"`
	LockedById                 *uint                 `json:"lockedById" gorm:"default:null"`
	LockedBy                   User                  `gorm:"foreignKey:LockedById" json:"lockedBy"`
	SelectedVehicles           []DriverVehicleAssign `gorm:"foreignKey:PlannerID" json:"selectedVehicles"`
	Name                       string                `json:"name"`
	Description                string                `json:"description"`
	DepotStartId               uint                  `json:"depotStartId" gorm:"default:null"`
	DepotEndId                 uint                  `json:"depotEndId" gorm:"default:null"`
	DepotStart                 Depot                 `gorm:"foreignKey:DepotStartId" json:"depotStart"`
	DepotEnd                   Depot                 `gorm:"foreignKey:DepotEndId" json:"depotEnd"`
	ServiceTime                uint                  `json:"serviceTime"`
	SetupTime                  uint                  `json:"setupTime"`
	Status                     string                `json:"status"` // [pending, optimized, published, in_progress, completed]
	TotalOrdersAssigned        uint                  `json:"totalOrdersAssigned"`
	TotalOrdersUnassigned      uint                  `json:"totalOrdersUnassigned"`
	TotalVehiclesAssigned      uint                  `json:"totalVehiclesAssigned"`
	TotalVehiclesUnassigned    uint                  `json:"totalVehiclesUnassigned"`
	TotalTime                  uint                  `json:"totalTime"`
	TotalDistance              uint                  `json:"totalDistance"`
	TotalWeight                uint                  `json:"totalWeight"`
	TotalWeightUnassigned      uint                  `json:"totalWeightUnassigned"`
	MinOrderWeight             uint                  `json:"minOrderWeight"`
	MaxOrderWeight             uint                  `json:"maxOrderWeight"`
	AverageWeightPerVehicle    float64               `json:"averageWeightPerVehicle"`
	AverageOrdersPerVehicle    float64               `json:"averageOrdersPerVehicle"`
	MinOrdersByVehicle         uint                  `json:"minOrdersByVehicle"`
	MaxOrdersByVehicle         uint                  `json:"maxOrdersByVehicle"`
	MinWeightByVehicle         uint                  `json:"minWeightByVehicle"`
	MaxWeightByVehicle         uint                  `json:"maxWeightByVehicle"`
	MinDistanceByVehicle       uint                  `json:"minDistanceByVehicle"`
	MaxDistanceByVehicle       uint                  `json:"maxDistanceByVehicle"`
	AdditionalVehiclesRequired float64               `json:"additionalVehiclesRequired"`
}
