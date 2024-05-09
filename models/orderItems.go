// Order Items
package models

import "gorm.io/gorm"

type OrderItem struct {
	*gorm.Model
	OrderID       uint    `json:"orderId"`                       // this is the foreign key for orders
	ItemID        uint    `json:"itemId"`                        // this is the foreign key for items
	Item          Item    `gorm:"foreignKey:ItemID" json:"item"` // this is the foreign key for items
	QuantityMain  int     `json:"quantityMain"`                  // this is the quantitymajor of the item
	QuantityMinor int     `json:"quantityMinor"`                 // this is the quantityminor of the item
	Capacity      int     `json:"totalCapacity"`                 // this is the total capacity of the item
	Weight        float64 `json:"weight"`                        // this is the weight of the item
	Volume        float64 `json:"volume"`                        // this is the volume of the item
	IsDelivered   bool    `json:"isDelivered"`                   // this is to check if the item is delivered
}
