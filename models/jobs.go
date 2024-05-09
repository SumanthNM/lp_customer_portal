/**
*
**/

package models

import (
	"time"

	"gorm.io/gorm"
)

type Jobs struct {
	*gorm.Model
	PlannerID             uint                `json:"plannerId"`
	JobId                 int64               `json:"jobId"`
	OrderID               *uint               `json:"orderId"`
	Order                 Order               `gorm:"foreignKey:OrderID" json:"order"`
	IsDriverAlocated      *bool               `json:"isDriverAllocated"`
	AllocatedDriver       uint                `json:"userID" gorm:"default:null"`
	Driver                User                `gorm:"foreignKey:AllocatedDriver;default:null" json:"user"`
	VehicleID             uint                `json:"vehicleId" gorm:"default:null"`
	Vehicle               Vehicle             `gorm:"foreignKey:VehicleID" json:"vehicle"`
	DriverVehicleAssignID *uint               `json:"driverVehicleAssignId" gorm:"default:null"`
	DriverVehicleAssign   DriverVehicleAssign `gorm:"foreignKey:DriverVehicleAssignID" json:"AllocatedDriver"`
	Sequence              int                 `json:"sequence"`    // defines the sequence of the job in the plan
	OldSequence           int                 `json:"oldSequence"` // defines the old sequence of the job in the plan [stores old seqeunce when the job is swapped] [used for audit trail]
	JobType               string              `json:"jobType"`     // whether pickup or delivery
	JobStatus             string              `json:"jobStatus"`   // whether pending, in progress, completed, cancelled
	SetupTime             int64               `json:"setupTime"`
	ServiceTime           time.Duration       `json:"serviceTime"`
	WaitingTime           int64               `json:"waitingTime"`
	EstimatedStartTime    time.Time           `json:"estimatedStartTime" gorm:"type:timestamp without time zone"`
	EstimatedEndTime      time.Time           `json:"estimatedEndTime" gorm:"type:timestamp without time zone"`
	ActualStartTime       int                 `json:"actualStartTime"`
	ActualEndTime         int                 `json:"actualEndTime"`
	EstimatedDuration     time.Duration       `json:"estimatedDuration"` // we calculate duration and store for easy queries
	ActualDuration        time.Duration       `json:"actualDuration"`
	IsManual              bool                `json:"isManual"` // determines whether the job is manually created or not
	IsCancelled           bool                `json:"isCancelled"`
	CancelledBy           uint                `json:"cancelledBy" gorm:"default:null"`             // contains the user id who cancelled the job
	CancelledUser         User                `gorm:"foreignKey:CancelledBy" json:"cancelledUser"` // contains the user details who cancelled the job
	CancelledReason       string              `json:"cancelledReason"`                             // contains the reason for cancellation
	PODSignature          []byte              `json:"podSignature"`                                // contains the signature of the customer
	PODImage              []byte              `json:"podImage"`                                    // contains the image of the delivery goods
	OTP                   string              `json:"otp"`                                         // contains the otp for the job
	IsOTPVerified         bool                `json:"isOtpVerified"`                               // determines whether the otp is verified or not
	Incidents             []Incident          `gorm:"foreignKey:JobID" json:"incidents"`           // contains the incidents for the job
	Lat                   float64             `json:"lat"`
	Lng                   float64             `json:"lng"`
	Geometry              string              `json:"geometry"`
}
