package customer_resource

import (
	"lp_customer_portal/common"
	"lp_customer_portal/schemas"
	"lp_customer_portal/services/customer_services"
	"net/http"
	"strconv"

	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

type CustomerResource struct {
}

func (cr *CustomerResource) CreateCustomer(ctx *restful.Context) {
	openlog.Info("Received request to create customer")
	var customerPayload schemas.CustomerPayload
	err := ctx.ReadEntity(&customerPayload)
	if err != nil {
		openlog.Error("Error occured while reading request body")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading request body"}, common.JSON_HEADER)
		return
	}
	// validate request body
	payloadErrs := common.ValidateStruct(customerPayload)
	if len(payloadErrs) > 0 {
		openlog.Error("Invalid request body")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Invalid request body", Data: payloadErrs}, common.JSON_HEADER)
		return
	}
	// insert customer
	cs := customer_services.New()
	res := cs.CreateCustomer(customerPayload)
	ctx.WriteHeaderAndJSON(res.Status, res, common.JSON_HEADER)
}

func (cr *CustomerResource) GetAllCustomers(ctx *restful.Context) {
	openlog.Info("Received request to get all customers")
	// get query params
	pageno, limit := common.GetPaginationParams(ctx.ReadQueryParameter("pageno"), ctx.ReadQueryParameter("limit"))
	filters := ctx.ReadQueryParameter("filters")
	// get customers
	cs := customer_services.New()
	res := cs.GetAllCustomers(pageno, limit, filters)
	ctx.WriteHeaderAndJSON(res.Status, res, common.JSON_HEADER)
}

func (cr *CustomerResource) GetCustomerById(ctx *restful.Context) {
	openlog.Info("Received request to get customer by id")
	// get id
	id := ctx.ReadPathParameter("id")
	cId, err := strconv.Atoi(id)
	if err != nil {
		openlog.Error("Error occured while converting id to integer")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Invalid Customer ID"}, common.JSON_HEADER)
		return
	}
	// get customer
	cs := customer_services.New()
	res := cs.GetCustomerById(cId)
	ctx.WriteHeaderAndJSON(res.Status, res, common.JSON_HEADER)
}

func (cr *CustomerResource) UpdateCustomer(ctx *restful.Context) {
	openlog.Info("Received request to update customer")
	// get id
	id := ctx.ReadPathParameter("id")
	cId, err := strconv.Atoi(id)
	if err != nil {
		openlog.Error("Invalid Customer ID")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Invalid Customer ID"}, common.JSON_HEADER)
		return
	}
	// get customer paylaod
	var customerPayload schemas.CustomerPayload
	err = ctx.ReadEntity(&customerPayload)
	if err != nil {
		openlog.Error("Error occured while reading request body")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading request body"}, common.JSON_HEADER)
		return
	}
	// validate request body
	payloadErrs := common.ValidateStruct(customerPayload)
	if len(payloadErrs) > 0 {
		openlog.Error("Invalid request body")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Invalid request body", Data: payloadErrs}, common.JSON_HEADER)
		return
	}
	// update customer
	cs := customer_services.New()
	res := cs.UpdateCustomerById(cId, customerPayload)
	ctx.WriteHeaderAndJSON(res.Status, res, common.JSON_HEADER)
}

func (cr *CustomerResource) DeleteCustomer(ctx *restful.Context) {
	openlog.Info("Received request to delete customer")
	// get id
	id := ctx.ReadPathParameter("id")
	cId, err := strconv.Atoi(id)
	if err != nil {
		openlog.Error("Invalid Customer ID")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Invalid Customer ID"}, common.JSON_HEADER)
		return
	}
	// delete customer
	cs := customer_services.New()
	res := cs.DeleteCustomerById(cId)
	ctx.WriteHeaderAndJSON(res.Status, res, common.JSON_HEADER)
}

func (cr *CustomerResource) Inject(service customer_services.CustomerServiceInterface) {

}

func (cr *CustomerResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodPost, Path: common.Basepath + "/customers", ResourceFunc: cr.CreateCustomer, Consumes: common.API_HEADERS, Produces: common.API_HEADERS, Read: schemas.CustomerPayload{}},
		{Method: http.MethodGet, Path: common.Basepath + "/customers", ResourceFunc: cr.GetAllCustomers, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: http.MethodGet, Path: common.Basepath + "/customers/{id}", ResourceFunc: cr.GetCustomerById, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: http.MethodPut, Path: common.Basepath + "/customers/{id}", ResourceFunc: cr.UpdateCustomer, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: http.MethodDelete, Path: common.Basepath + "/customers/{id}", ResourceFunc: cr.DeleteCustomer, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
	}
}
