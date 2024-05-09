package schemas

type SelectAllOrdersPayload struct {
	IsSelected         uint   `json:"isSelected" validate:"required"`
	PreferredStartDate string `json:"preferredStartDate" validate:"required"`
	PreferredEndDate   string `json:"preferredEndDate" validate:"required"`
}
