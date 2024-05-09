package schemas

type OrderStatusPayload struct {
	PreferredDateTime    string `json:"preferredDateTime" validate:"required"`
	PreferredEndDateTime string `json:"preferredEndDateTime" validate:"required"`
}
