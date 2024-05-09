/**
* the output for the jobs for time line
**/

package models

type VehicleTable struct {
	DriverVehicleAssignID uint
	LicensePlate          string
	DriverName            string
	NoOfOrders            uint
	NoOfJobs              uint
	NoOfSKU               uint
	TotalWeight           uint
	OptimizationRequired  bool
}
