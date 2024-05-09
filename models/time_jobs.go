/**
* the output for the jobs for time line
**/

package models

import "time"

type TimeJobs struct {
	DriverID                   int       `json:"driverId"`
	JobID                      int       `json:"jobId"`
	JobType                    string    `json:"jobType"`
	OrderID                    string    `json:"orderId"`
	VehicleId                  int       `json:"vehicleId"`
	EstimatedStartTime         time.Time `json:"estimatedStartTime"`
	EstimatedEndTime           time.Time `json:"estimatedEndTime"`
	EstimatedDuration          int64     `json:"estimatedDuration"`
	AddressStr                 string    `json:"addressStr"`
	Lat                        float64   `json:"lat"`
	Lng                        float64   `json:"lng"`
	Color                      string    `json:"color"`
	SequenceNo                 int       `json:"sequence"`
	CustomerName               string    `json:"customerName"`
	JobStatus                  string    `json:"jobStatus"`
	OrderNo                    string    `json:"orderNo"`
	Geometry                   string    `json:"geometry"`
	ServiceTime                int64     `json:"serviceTime"`
	PreferredDeliveryStartDate time.Time `json:"preferredDeliveryStartDate"` // maintains the date and time for start
	PreferredDeliveryEndDate   time.Time `json:"preferredDeliveryEndDate"`   // maintains the date and time for end
	PreferredPickupStartDate   time.Time `json:"preferredPickupStartDate"`   // maintains the date and time for start
	PreferredPickupEndDate     time.Time `json:"preferredPickupEndDate"`     // maintains the date and time for end
	IsManual                   bool      `json:"isManual"`
}

// used for filtering TimeJobs
type DateTimeRange struct {
	EstimatedStartTime string `json:"estimatedStartTime"`
	EstimatedEndTime   string `json:"estimatedEndTime"`
}