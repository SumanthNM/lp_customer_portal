/**
* the output for the jobs for time line
**/

package models

import "time"

type UnassignedJobs struct {
	JobID                      int       `json:"jobId"`
	JobType                    string    `json:"jobType"`
	OrderID                    string    `json:"orderId"`
	AddressStr                 string    `json:"addressStr"`
	Lat                        string    `json:"lat"`
	Lng                        string    `json:"lng"`
	CustomerName               string    `json:"customerName"`
	OrderNo                    string    `json:"orderNo"`
	IsManual                   bool      `json:"isManual"`
	ServiceTime                int64     `json:"serviceTime"`
	PreferredDeliveryStartDate time.Time `json:"preferredDeliveryStartDate"` // maintains the date and time for start
	PreferredDeliveryEndDate   time.Time `json:"preferredDeliveryEndDate"`   // maintains the date and time for end
	PreferredPickupStartDate   time.Time `json:"preferredPickupStartDate"`   // maintains the date and time for start
	PreferredPickupEndDate     time.Time `json:"preferredPickupEndDate"`     // maintains the date and time for end
}
