/**
 * sample user model
 *  All DB models should be defined in folder models
**/

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model                              // <-- this will add id, created_at, updated_at, deleted_at
	Name              string                 `json:"name"`
	Age               int                    `json:"age"`
	Email             string                 `json:"email"`
	Password          string                 `json:"password"`
	Salt              string                 `json:"salt"`
	Contact           string                 `json:"contact"`
	ContractStartDate time.Time              `json:"contractStartDate"`
	ContractEndDate   time.Time              `json:"contractEndDate"`
	VerificationCode  string                 `json:"verificationCode"`
	IsActive          bool                   `json:"isActive"`
	OnDuty            bool                   `json:"onDuty"`
	ProfilePicture    []byte                 `json:"profilePicture"`
	Skillsets         []*Skillset            `gorm:"many2many:user_skillsets;" json:"skillsets"`
	Roles             []*Role                `gorm:"many2many:user_roles;" json:"roles"`
	Zone              []*Zone                `gorm:"many2many:user_zones" json:"zones"`
	AssignedVehicles  []*DriverVehicleAssign `gorm:"foreignKey:UserID" json:"vehicles"`
	Shifts            []*Shift               `gorm:"many2many:user_shifts" json:"shifts"`
	// temp field for versafleet
	VfID              uint                   `json:"vfId"`
}
