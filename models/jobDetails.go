/**
*
**/

package models

import (
	"time"
)

type JobDetails struct {
	VehicleID             uint
	DriverVehicleAssignID uint
	DriverName            string
	LicensePlate          string
	OrderNo               string
	CustomerName          string
	JobType               string
	Address               string
	PostalCode            int
	PreferredStartTime    time.Time
	PreferredEndTime      time.Time
	EstimatedStartTime    time.Time
	EstimatedEndTime      time.Time
	Item                  string
	Unit                  string
	Quantity              string
	Remarks               string
}
