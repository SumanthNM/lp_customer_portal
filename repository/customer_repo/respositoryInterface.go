/**
 * Customer Repository Interface
 *
**/

package customer_repo

import "lp_customer_portal/models"

type CustomerRepositoryInterface interface {
	Insert(customer models.Customer) (models.Customer, error)
	GetAllCustomers(pageno, limit int, filters string) ([]models.Customer, error)
	GetCustomerByEmail(email string) (bool, error)
	GetCustomerById(id int) (models.Customer, error)
	UpdateCustomerById(id int, customer models.Customer) (models.Customer, error)
	DeleteCustomerById(id int) error
	GetCount() (int64, error)
	InsertAddress(address models.Address) (models.Address, error)
	GetAllAddress(pageno, limit int) ([]models.Address, error)
	BulkInsertAddresses(addresses []models.Address) (int64, error)
	FetchUniqueAddress(query string) ([]models.Address, error)
	BulkInsertCustomer(customers []models.Customer) (int64, error)
	GetAddressByStreetAddress(addr string) (models.Address, error)
}
