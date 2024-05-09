package user_repo

import (
	"lp_customer_portal/common"
	"lp_customer_portal/models"

	"github.com/go-chassis/openlog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	DB *gorm.DB
}

// Fetch Selected Users
func (ur UserRepo) FetchSelectedUsers(users []int) ([]models.User, error) {
	openlog.Debug("Fetching selected users from database")
	var usersList []models.User
	res := ur.DB.Preload(clause.Associations).Model(&models.User{}).Where("id IN ?", users).Find(&usersList)
	if res.Error != nil {
		openlog.Error("Error while fetching selected users " + res.Error.Error())
		return usersList, res.Error
	}
	return usersList, nil
}

// Fetch Selected User
func (ur UserRepo) GetUserById(userId int) (models.User, error) {
	openlog.Info("Fetching user by id from database")
	var user models.User
	result := ur.DB.Model(&models.User{}).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			openlog.Error("User not found")
			return user, common.ErrResourceNotFound
		}
		openlog.Error("Error occured while fetching user by id from database")
		return user, result.Error
	}
	return user, nil
}
