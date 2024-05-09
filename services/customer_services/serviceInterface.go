/**
 * Customer service interface
 *
 *
**/
package customer_services

import (
	"lp_customer_portal/common"
	"lp_customer_portal/schemas"
)

type CustomerServiceInterface interface {
	CreateCustomer(customer schemas.CustomerPayload) common.HTTPResponse
	GetAllCustomers(pageno, limit int, filters string) common.HTTPResponse
	GetCustomerById(id int) common.HTTPResponse
	UpdateCustomerById(id int, customer schemas.CustomerPayload) common.HTTPResponse
	DeleteCustomerById(id int) common.HTTPResponse
	CreateAddress(address schemas.Address) common.HTTPResponse
}
