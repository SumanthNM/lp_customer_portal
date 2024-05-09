package models

import (
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	*gorm.Model

	ShiftName      string    `json:"shiftName"`
	ShiftStartTime time.Time `json:"shiftStartTime"`
	ShiftEndTime   time.Time `json:"shiftEndTime"`
	BreakStartTime time.Time `json:"breakStartTime"`
	BreakEndTime   time.Time `json:"breakEndTime"`
	BreakDuration  int       `json:"breakDuration"`
	Users          []*User   `gorm:"many2many:user_shifts;" json:"users"`
}
