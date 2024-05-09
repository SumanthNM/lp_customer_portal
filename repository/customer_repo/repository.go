/**
 * Customer Repo Implementation
 *
**/
package customer_repo

import (
	"lp_customer_portal/common"
	"lp_customer_portal/models"

	"github.com/go-chassis/openlog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerRepository struct {
	DB *gorm.DB
}

func (cr CustomerRepository) Insert(customer models.Customer) (models.Customer, error) {
	openlog.Debug("Inserting customer into database")
	res := cr.DB.Create(&customer)
	if res.Error != nil {
		openlog.Error("Error occured while inserting customer into database")
		return customer, res.Error
	}
	return customer, nil
}

func (cr CustomerRepository) GetAllCustomers(pageno, limit int, filters string) ([]models.Customer, error) {
	openlog.Debug("Fetching all customers from database")
	//TODO: add filters
	var customers []models.Customer
	filterScope := common.GetAllCondition(filters)
	result := cr.DB.Scopes(filterScope...).Model(&customers).Offset((pageno - 1) * limit).Limit(limit).Find(&customers) // adding pagination to query
	if result.Error != nil {
		openlog.Error("Error occured while fetching all customers from database")
		return customers, result.Error
	}
	return customers, nil
}

func (cr CustomerRepository) GetCustomerById(id int) (models.Customer, error) {
	openlog.Debug("Fetching customer from database")
	var customer models.Customer
	result := cr.DB.Where("id = ?", id).First(&customer)
	if result.Error != nil {
		openlog.Error("Error occured while fetching customer from database")
		return customer, result.Error
	}
	if result.RowsAffected == 0 {
		openlog.Error("customer not found")
		return customer, common.ErrResourceNotFound
	}

	return customer, nil
}

func (cr CustomerRepository) UpdateCustomerById(id int, customer models.Customer) (models.Customer, error) {
	openlog.Debug("Updating customer in database")
	res := cr.DB.Model(&customer).Where("id = ?", id).Updates(customer)
	if res.Error != nil {
		openlog.Error("Error occured while updating customer in database")
		if res.Error == gorm.ErrRecordNotFound {
			return customer, common.ErrResourceNotFound
		}
		return customer, res.Error
	}
	result := cr.DB.First(&customer, id)
	if result.Error != nil {
		openlog.Error("Error occurred while reloading customer")
		return customer, result.Error
	}
	return customer, nil
}

func (cr CustomerRepository) DeleteCustomerById(id int) error {
	openlog.Debug("Deleting customer from database")
	result := cr.DB.Where("id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		openlog.Error("Error occured while deleting customer from database")
		return result.Error
	}
	return nil
}

func (cr CustomerRepository) GetCount() (int64, error) {
	openlog.Debug("Fetching count of customers from database")
	var count int64
	res := cr.DB.Model(&models.Customer{}).Count(&count)
	if res.Error != nil {
		openlog.Error("Error occured while fetching count of customers from database")
		return count, res.Error
	}
	return count, nil
}

func (cr CustomerRepository) GetCustomerByEmail(email string) (bool, error) {
	openlog.Debug("Fetching customer from database")
	var customer models.Customer
	res := cr.DB.Where("email = ?", email).First(&customer)
	if res.Error != nil {
		openlog.Error("Error occured while fetching customer from database")
		if res.Error == gorm.ErrRecordNotFound {
			return false, common.ErrResourceNotFound
		}
	}
	return false, nil
}

func (cr CustomerRepository) InsertAddress(address models.Address) (models.Address, error) {
	openlog.Debug("Inserting address into database")
	err := cr.DB.Create(&address).Error
	return address, err
}

// get all address from database
func (cr CustomerRepository) GetAllAddress(pageno, limit int) ([]models.Address, error) {
	openlog.Info("Fetching all address from database")
	var address []models.Address
	result := cr.DB
	if pageno > 0 {
		result = result.Offset((pageno - 1) * limit).Limit(limit)
	}

	result = result.Find(&address) // adding pagination to query
	if result.Error != nil {
		openlog.Error("Error occured while fetching all address from database")
		return address, result.Error
	}
	return address, nil
}

func (cr CustomerRepository) BulkInsertAddresses(addresses []models.Address) (int64, error) {
	openlog.Debug("Bulk inserting addresses")
	chunkSize := 100
	for i := 0; i < len(addresses); i += chunkSize {
		end := i + chunkSize
		if end > len(addresses) {
			end = len(addresses)
		}
		res := cr.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(addresses[i:end])
		if res.Error != nil {
			openlog.Error("Error occured while inserting the addresses")
			return 0, res.Error
		}
	}
	return int64(len(addresses)), nil
}

func (cr CustomerRepository) FetchUniqueAddress(query string) ([]models.Address, error) {
	var addressList []string

	res := cr.DB.Raw(query).Find(&addressList)
	if res.Error != nil {
		openlog.Error("Error while fetching selected addresses " + res.Error.Error())
		return []models.Address{}, res.Error
	}
	var addressModels []models.Address
	for _, address := range addressList {
		addressModel := models.Address{
			Address_Str: address,
		}
		addressModels = append(addressModels, addressModel)
	}

	return addressModels, nil

}

func (cr CustomerRepository) BulkInsertCustomer(customers []models.Customer) (int64, error) {
	openlog.Debug("Bulk inserting customers")
	chunkSize := 100
	for i := 0; i < len(customers); i += chunkSize {
		end := i + chunkSize
		if end > len(customers) {
			end = len(customers)
		}
		res := cr.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "customer_no"}}, // to ref unique sku column
			DoNothing: true,
		}).Create(customers[i:end])
		if res.Error != nil {
			openlog.Error("Error occured while inserting the customers")
			return 0, res.Error
		}
	}
	return int64(len(customers)), nil
}

func (cr CustomerRepository) GetAddressByStreetAddress(addr string) (models.Address, error) {
	openlog.Debug("Fetching customer from database")
	var address models.Address
	result := cr.DB.Where("address_str = ?", addr).First(&address)
	if result.RowsAffected == 0 {
		openlog.Error("address not found")
		return address, common.ErrResourceNotFound
	}
	if result.Error != nil {
		openlog.Error("Error occured while fetching address from database")
		return address, result.Error
	}
	return address, nil
}
