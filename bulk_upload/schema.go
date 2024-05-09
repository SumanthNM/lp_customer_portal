package bulk_upload

import "time"

type ItemRow struct {
	ItemSKU   string  `json:"itemSKU"`
	Weight    float64 `json:"weight"`
	Volume    float64 `json:"volume"`
	Name      string  `json:"name"`
	InnerPack int     `json:"innerPack"`
}

type VehicleRow struct {
	LicensePlate string `json:"licensePlate"`
	VehicleType  int    `json:"vehicleType"`
	DepotStartId uint   `json:"depotStartId"`
	DepotEndId   uint   `json:"depotEndId"`
}

type OrderRow struct {
	OrderNo       string `json:"orderNo"`
	CustomerNo    string `json:"customerNo"`
	CustomerName  string `json:"customerName"`
	CustomerEmail string `json:"customer_email"`
	QuantityMain  int    `json:"quantityMain"`
	// QuantityMinor int       `json:"quantityMinor"`
	OrderedDate time.Time `json:"date"`
	// PreferredDeliverStartTime time.Time `json:"preferredDeliverStartTime"`
	// PreferredDeliverEndTime   time.Time `json:"preferredDeliverEndTime"`
	Priority string `json:"priority"`
	// Weight                    float64   `json:"weight"`
	// Volume                    float64   `json:"volume"`
	// ZipCode                   string    `json:"zipCode"`
	// AddressStr                string    `json:"addressStr"`
	ItemRow                        ItemRow   `json:"itemRow"`
	PickupAddress                  string    `json:"pickupAddress"`
	PickupLat                      float64   `json:"pickupLat"`
	PickupLng                      float64   `json:"pickupLng"`
	PickupPostal                   int       `json:"pickupPostal"`
	PreferredPickupStartDateTime   time.Time `json:"preferredPickupStartTime"`
	PreferredPickupEndDateTime     time.Time `json:"preferredPickupEndTime"`
	PickupCustomServiceTime        int       `json:"pickupCustomServiceTime"`
	DeliveryAddress                string    `json:"deliveryAddress"`
	DeliveryPostal                 int       `json:"deliveryPostal"`
	DeliveryLat                    float64   `json:"deliveryLat"`
	DeliveryLng                    float64   `json:"deliveryLng"`
	PreferredDeliveryStartDateTime time.Time `json:"preferredDeliveryStartTime"`
	PreferredDeliveryEndDateTime   time.Time `json:"preferredDeliveryEndTime"`
	DeliveryCustomServiceTime      int       `json:"deliveryCustomServiceTime"`
	Zone                           int       `json:"zone"`
	Skills                         int       `json:"skills"`
	JobType                        string    `json:"jobType"`
}
