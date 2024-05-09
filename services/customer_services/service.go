/**
 * Customer service Implementation
 *
**/
package customer_services

import (
	"lp_customer_portal/common"
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	"lp_customer_portal/repository/customer_repo"
	"lp_customer_portal/schemas"

	"github.com/go-chassis/openlog"
)

type CustomerService struct {
	Repo customer_repo.CustomerRepositoryInterface
}

func New() *CustomerService {
	openlog.Info("Initializing Customer Service")
	db := database.GetClient()
	CustomerRepo := customer_repo.CustomerRepository{DB: db}
	return &CustomerService{
		Repo: &CustomerRepo,
	}
}

func (cs *CustomerService) CreateCustomer(customerPayload schemas.CustomerPayload) common.HTTPResponse {
	openlog.Debug("Inserting customer")
	// check for duplicates
	exists, err := cs.Repo.GetCustomerByEmail(customerPayload.Email)
	if err != nil && err != common.ErrResourceNotFound { // if error is not "resource not found" then return error
		openlog.Error("Error occured while checking for duplicates")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while checking for duplicates"}
	}
	if exists {
		openlog.Error("Duplicate customer found")
		return common.HTTPResponse{Status: 409, Msg: "Duplicate customer found"}
	}
	// convert schema to model
	customerModel := models.Customer{}
	err = customerPayload.ToModel(&customerModel)
	if err != nil {
		openlog.Error("Error occured while converting payload to model")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while converting payload to model"}
	}
	// insert into database
	customer, err := cs.Repo.Insert(customerModel)
	if err != nil {
		openlog.Error("Error occured while inserting customer into database")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while inserting customer into database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Customer inserted successfully", Data: customer}
}

func (cs *CustomerService) GetAllCustomers(pageno, limit int, filters string) common.HTTPResponse {
	openlog.Debug("Fetching all customers")
	// get customers from database
	customers, err := cs.Repo.GetAllCustomers(pageno, limit, filters)
	if err != nil {
		openlog.Error("Error occured while fetching customers from database")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching customers from database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Customers fetched successfully", Data: customers}
}

func (cs *CustomerService) GetCustomerById(id int) common.HTTPResponse {
	openlog.Debug("Fetching customer by id")
	// get customer from database
	customer, err := cs.Repo.GetCustomerById(id)
	if err != nil {
		openlog.Error("Error occured while fetching customer from database")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Customer not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching customer from database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Customer fetched successfully", Data: customer}
}

func (cs *CustomerService) UpdateCustomerById(id int, customerPayload schemas.CustomerPayload) common.HTTPResponse {
	openlog.Debug("Updating customer by id")
	// check if customer exists
	_, err := cs.Repo.GetCustomerById(id)
	if err != nil {
		openlog.Error("Error occured while fetching customer from database")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Customer not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching customer from database"}
	}
	// convert schema to model
	customerModel := models.Customer{}
	err = customerPayload.ToModel(&customerModel)
	if err != nil {
		openlog.Error("Error occured while converting payload to model")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while converting payload to model"}
	}
	// update customer in database
	customer, err := cs.Repo.UpdateCustomerById(id, customerModel)
	if err != nil {
		openlog.Error("Error occured while updating customer in database")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Customer not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while updating customer in database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Customer updated successfully", Data: customer}
}

func (cs *CustomerService) DeleteCustomerById(id int) common.HTTPResponse {
	openlog.Debug("Deleting customer by id")
	// check if customer exists
	_, err := cs.Repo.GetCustomerById(id)
	if err != nil {
		openlog.Error("Error occured while fetching customer from database")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Customer not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching customer from database"}
	}
	// delete customer from database
	err = cs.Repo.DeleteCustomerById(id)
	if err != nil {
		openlog.Error("Error occured while deleting customer from database")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while deleting customer from database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Customer deleted successfully"}
}

func (cs *CustomerService) CreateAddress(addressPayload schemas.Address) common.HTTPResponse {
	openlog.Debug("Inserting address")
	// convert schema to model
	addressModel := models.Address{}
	err := addressPayload.ToModel(&addressModel)
	if err != nil {
		openlog.Error("Error occured while converting payload to model")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while converting payload to model"}
	}
	// TODO: check if address lat and long is given.
	// if lat long is given add zone to the addaddressModel
	// TODO: if lat and long is not given, then user here api to get lat and long from address
	// insert into database
	address, err := cs.Repo.InsertAddress(addressModel)
	if err != nil {
		openlog.Error("Error occured while inserting address into database")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while inserting address into database"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Address inserted successfully", Data: address}
}
