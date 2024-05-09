// Incidents Models

package models

import (
	"gorm.io/gorm"
)

type Incident struct {
	*gorm.Model
	Type        string  `json:"type"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	ReportedBy  uint    `json:"reportedBy"`
	Reported    User    `gorm:"foreignKey:ReportedBy" json:"reported"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
	JobID       int     `json:"jobId"`
	Jobs        Jobs    `gorm:"foreignKey:JobID" json:"job"`
}
