/**
*
*
**/

package models

import (
	//"lp_oms/models"
	"time"

	"gorm.io/gorm"
)

// Order is the db model for orders table
type Order struct {
	*gorm.Model
	// OrderNo    string      `json:"orderNo" gorm:"unique"` // order_no is the job_id for delivery, needed for pyvroom // unique.

	// temp: make orderno not unique to allow for push into VF db
	OrderNo string `json:"orderNo"` // order_no is the job_id for delivery, needed for pyvroom // unique.

	ContractNo string      `json:"contractNo"`
	OrderDate  time.Time   `json:"orderDate"`  // new field
	OrderType  string      `json:"orderType"`  // new field [pickup or delivery or pickup and delivery]
	OrderItems []OrderItem `json:"orderItems"` // this is the foreign key , | Capacity of orders are from the Items

	PickupAddressID uint     `json:"pickupAddressId" gorm:"default:null"`
	PickupAddress   *Address `json:"pickupAddress" gorm:"foreignKey:PickupAddressID"`

	DeliveryAddressID uint     `json:"deliveryAddressId" gorm:"default:null"`
	DeliveryAddress   *Address `json:"deliveryAddress" gorm:"foreignKey:DeliveryAddressID"`

	Zone   *Zone `json:"zone" gorm:"foreignKey:ZoneID"`
	ZoneID uint  `json:"zoneId" gorm:"default:null"`
	// Customer and Address Information.
	CustomerID uint     `json:"customerId"` // this is the foreign key
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID"`

	//AddressID uint     `json:"addressId" gorm:"default:null"` // this is the foreign key
	//Address   *Address `gorm:"foreignKey:AddressID" json:"address"`

	// preferred date and time for  both start and end
	PreferredDeliveryStartDate time.Time `json:"preferredDeliveryStartDate" ` // maintains the date and time for start
	PreferredDeliveryEndDate   time.Time `json:"preferredDeliveryEndDate" `   // maintains the date and time for end
	DeliveryServiceTime        int       `json:"deliveryServiceTime"`
	PreferredPickupStartDate   time.Time `json:"preferredPickupStartDate" ` // maintains the date and time for start
	PreferredPickupEndDate     time.Time `json:"preferredPickupEndDate" `   // maintains the date and time for end
	PickupServiceTime          int       `json:"pickupServiceTime"`

	// this is technician who place the order
	StaffID uint  `json:"staffId" gorm:"default:null"` // this is the foreign key
	Staff   *User `json:"assignedStaff" gorm:"foreignKey:StaffID"`

	// this is to record the driver who will fulfil the trip
	DriverVehicleAssignID uint                 `json:"AssignedVehicle" gorm:"default:null"` // this is the foreign key
	DriverVehicleAssign   *DriverVehicleAssign `gorm:"foreignKey:DriverVehicleAssignID" json:"DriverVehicleAssign"`

	// skillset is not relevant for Kone technical
	SkillsetID uint      `json:"skillsetId" gorm:"default:null"` // this is the foreign key
	Skillset   *Skillset `json:"skillset"`                       // this is the foreign key

	// the priority of the trip, Kone will have different SLA for each priority
	Priority string `json:"priority"`

	// have this trip schedule for a driver to fulfil it
	IsScheduled bool `json:"isScheduled"`

	// is the technician in transit
	InTransit bool `json:"inTransit"`

	// is this order completed (drop off for Kone)
	IsCompleted bool `json:"isCompleted"`
	// is this order selected for the optimization.
	PlannerID  uint    `json:"plannerId" gorm:"default:null"` // this is the foreign key
	IsSelected Planner `json:"isSelected" gorm:"foreignKey:PlannerID"`

	// SLA calculation
	ActualTravelTimeSecond int     `json:"actualTravelTimeSecond"`
	ActualTimeSinceCall    int     `json:"actualTimeSinceCall"`
	SLAPassFail            bool    `json:"SLAPassFail"`
	Archived               bool    `json:"archived"` // is this order archived
	Status                 string  `json:"status"`
	Volume                 float64 `json:"volume"`
	Weight                 float64 `json:"weight"`

	// if user's actions are req to deduplicate orders
	Duplicates bool `json:"duplicates"`
	// temp field to for versafleet
	VfID    uint   `json:"vfId"`
	Remarks string `json:"remarks"`
}
