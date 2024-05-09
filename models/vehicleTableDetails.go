/**
* the output for the jobs for time line
**/

package models

type VehicleTableDetails struct {
	OrderNo          string
	ItemName         string
	SKU              string
	OrderItemWeight  uint
	TotalOrderWeight uint
	QuantityMain     uint
	QuantityMinor    uint
	InnerPack        uint
}
