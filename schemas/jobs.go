package schemas

import (
	"errors"
	"lp_customer_portal/common"
	"lp_customer_portal/models"
	"time"
)

type Jobs struct {
	PlannerID uint `json:"plannerId"`
	//JobId                 int64   `json:"jobId"`
	OrderID               *uint   `json:"orderId"`
	IsDriverAlocated      bool    `json:"isDriverAllocated"`
	DriverVehicleAssignID *uint   `json:"driverVehicleAssignId"`
	VehicleID             uint    `json:"vehicleId"`
	Sequence              int     `json:"sequence"`    // defines the sequence of the job in the plan
	OldSequence           int     `json:"oldSequence"` // defines the old sequence of the job in the plan [stores old seqeunce when the job is swapped] [used for audit trail]
	JobType               string  `json:"jobType"`     // whether pickup or delivery
	JobStatus             string  `json:"jobStatus"`   // whether pending, in progress, completed, cancelled
	SetupTime             int64   `json:"setupTime"`
	ServiceTime           int64   `json:"serviceTime"`
	WaitingTime           int64   `json:"waitingTime"`
	EstimatedStartTime    string  `json:"estimatedStartTime"`
	EstimatedEndTime      string  `json:"estimatedEndTime"`
	ActualStartTime       int     `json:"actualStartTime"`
	ActualEndTime         int     `json:"actualEndTime"`
	EstimatedDuration     int64   `json:"estimatedDuration"` // we calculate duration and store for easy queries
	ActualDuration        int64   `json:"actualDuration"`
	IsManual              bool    `json:"isManual"` // determines whether the job is manually created or not
	IsCancelled           bool    `json:"isCancelled"`
	CancelledBy           uint    `json:"cancelledBy"`     // contains the user id who cancelled the job
	CancelledReason       string  `json:"cancelledReason"` // contains the reason for cancellation
	PODSignature          []byte  `json:"podSignature"`    // contains the signature of the customer
	PODImage              []byte  `json:"podImage"`        // contains the image of the delivery goods
	OTP                   string  `json:"otp"`             // contains the otp for the job
	IsOTPVerified         bool    `json:"isOtpVerified"`   // determines whether the otp is verified or not
	Lat                   float64 `json:"lat"`
	Lng                   float64 `json:"lng"`
}

func (j *Jobs) ToModel(model *models.Jobs) error {

	if j.EstimatedStartTime != "" {
		estimatedStartTime, e := time.Parse(common.TIME_FORMAT, j.EstimatedStartTime)
		if e != nil {
			return errors.New("error while parsing estimated delivery time")
		}
		model.EstimatedStartTime = estimatedStartTime
	}

	if j.EstimatedEndTime != "" {
		estimatedEndTime, e := time.Parse(common.TIME_FORMAT, j.EstimatedEndTime)
		if e != nil {
			return errors.New("error while parsing estimated pickup time")
		}
		model.EstimatedEndTime = estimatedEndTime
	}

	model.EstimatedDuration = time.Duration(j.EstimatedDuration)
	model.ActualDuration = time.Duration(j.ActualDuration)
	model.ServiceTime = time.Duration(j.ServiceTime)
	model.IsDriverAlocated = &j.IsDriverAlocated
	model.DriverVehicleAssignID = j.DriverVehicleAssignID
	model.Sequence = j.Sequence
	model.OldSequence = j.OldSequence
	model.JobType = j.JobType
	model.JobStatus = j.JobStatus
	model.SetupTime = j.SetupTime
	model.WaitingTime = j.WaitingTime
	model.ActualStartTime = j.ActualStartTime
	model.ActualEndTime = j.ActualEndTime
	model.IsCancelled = j.IsCancelled
	model.CancelledBy = j.CancelledBy
	model.CancelledReason = j.CancelledReason
	model.PODSignature = j.PODSignature
	model.PODImage = j.PODImage
	model.OTP = j.OTP
	model.IsOTPVerified = j.IsOTPVerified
	model.Lat = j.Lat
	model.Lng = j.Lng
	return nil
}
