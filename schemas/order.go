package schemas

import (
	"fmt"
	"lp_customer_portal/common"
	"lp_customer_portal/models"
	"time"
)

type OrderItemPayload struct {
	ItemID        int     `json:"itemId"`
	QuantityMain  int     `json:"quantityMain"`
	QuantityMinor int     `json:"quantityMinor"`
	Capacity      int     `json:"totalCapacity"`
	Weight        float64 `json:"weight"`
	Volume        float64 `json:"volume"`
}

type OrderPayload struct {
	OrderNo                    string             `json:"orderNo"`
	ContractNo                 string             `json:"contractNo"`
	OrderItems                 []OrderItemPayload `json:"orderItems"`
	CustomerID                 uint               `json:"customerId"` // this will be null if the customer is new.
	Customer                   CustomerPayload    `json:"customer" `
	AddressID                  uint               `json:"addressId"` // this will be null if the Address is new, if address is new, customer will be added with the new address.
	Address                    Address            `json:"address"`
	PreferredDeliveryStartDate string             `json:"preferredDeliveryStartDate" ` // maintains the date and time for start
	PreferredDeliveryEndDate   string             `json:"preferredDeliveryEndDate" `
	PreferredPickupStartDate   string             `json:"preferredPickupStartDate"` // maintains the date and time for start
	PreferredPickupEndDate     string             `json:"preferredPickupEndDate"`   // maintains the date and time for end
	SkillsetID                 uint               `json:"skillsetId"`
	StaffID                    uint               `json:"staffId"`
	Priority                   string             `json:"priority"` // dropdown urgent, normal, low
	DriverVehicleAssignID      uint               `json:"driverVehicleAssignId"`
	OrderDate                  string             `json:"orderDate"`
	Status                     string             `json:"status"`
}

func (op *OrderPayload) ToModel(orderModel *models.Order) error {
	orderModel.OrderNo = op.OrderNo
	orderModel.ContractNo = op.ContractNo
	orderModel.CustomerID = op.CustomerID
	orderModel.DeliveryAddressID = op.AddressID
	orderModel.Status = op.Status
	fmt.Println("*********", op.PreferredDeliveryStartDate)
	if op.PreferredDeliveryStartDate != "" {
		date, err := time.Parse(common.TIME_FORMAT, op.PreferredDeliveryStartDate)
		if err != nil {
			return err
		}
		fmt.Println(date)
		orderModel.PreferredDeliveryStartDate = date
	}
	if op.PreferredDeliveryEndDate != "" {
		date, err := time.Parse(common.TIME_FORMAT, op.PreferredDeliveryEndDate)
		if err != nil {
			return err
		}
		orderModel.PreferredDeliveryEndDate = date
	}
	if op.PreferredPickupStartDate != "" {
		date, err := time.Parse(common.TIME_FORMAT, op.PreferredPickupStartDate)
		if err != nil {
			return err
		}
		fmt.Println(date)
		orderModel.PreferredPickupStartDate = date
	}
	if op.PreferredPickupEndDate != "" {
		date, err := time.Parse(common.TIME_FORMAT, op.PreferredPickupEndDate)
		if err != nil {
			return err
		}
		orderModel.PreferredPickupEndDate = date
	}

	if op.OrderDate != "" {
		date, err := time.Parse(common.ORDER_TIME_FORMAT, op.OrderDate)
		if err != nil {
			return err
		}
		orderModel.OrderDate = date
	}
	orderModel.StaffID = op.StaffID
	orderModel.SkillsetID = op.SkillsetID //TODO Temporary
	orderModel.Priority = op.Priority
	orderModel.DriverVehicleAssignID = op.DriverVehicleAssignID
	if orderModel.OrderItems == nil {
		orderModel.OrderItems = make([]models.OrderItem, len(op.OrderItems))
	}
	for i, item := range op.OrderItems {
		orderModel.OrderItems[i].ItemID = uint(item.ItemID)
		orderModel.OrderItems[i].QuantityMain = item.QuantityMain
		orderModel.OrderItems[i].QuantityMinor = item.QuantityMinor
		orderModel.OrderItems[i].Capacity = item.Capacity
		orderModel.OrderItems[i].Weight = item.Weight
		orderModel.OrderItems[i].Volume = item.Volume
	}
	return nil
}
