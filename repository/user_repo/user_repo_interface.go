package user_repo

import "lp_customer_portal/models"

type UserRepoInterface interface {
	FetchSelectedUsers(users []int) ([]models.User, error)
	GetUserById(userId int) (models.User, error)
}
