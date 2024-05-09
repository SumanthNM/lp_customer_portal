package schemas

import (
	"lp_customer_portal/models"
)

type ItemPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category" validate:"required"`
	SKU         string  `json:"sku"`
	Capacity    float64 `json:"capacity"`
	Volume      float64 `json:"volume"`
	Weight      float64 `json:"weight"`
	Status      string  `json:"status"`
	InnerPack   int     `json:"innerPack"`
}

func (i *ItemPayload) ToModel(item *models.Item) error {
	item.Name = i.Name
	item.Description = i.Description
	item.Category = i.Category
	item.SKU = i.SKU
	item.Volume = i.Volume
	item.Weight = i.Weight
	item.InnerPack = i.InnerPack
	item.Status = i.Status

	return nil
}
