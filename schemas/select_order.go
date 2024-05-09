package schemas

type SelectOrder struct {
	IsSelected uint `json:"isSelected" validate:"required"`
}
