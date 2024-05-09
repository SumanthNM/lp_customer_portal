package order_resource

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"lp_customer_portal/bulk_upload"
	"lp_customer_portal/common"
	"lp_customer_portal/schemas"
	"lp_customer_portal/services/order_service"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

type OrderResource struct {
}

func (or *OrderResource) Inject(orderservice order_service.OrderServiceInterface) {

}

// create order
func (or *OrderResource) CreateOrder(ctx *restful.Context) {
	openlog.Info("Got request to create order")
	// read payload
	order := schemas.OrderPayload{}
	ctx.ReadEntity(&order)
	// validate paylaod
	payloadErrs := common.ValidateStruct(order)
	if len(payloadErrs) > 0 {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Payload validation error", Data: payloadErrs}, "application/json")
		return
	}
	// call service
	os := order_service.New()
	res := os.CreateOrder(order)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

// get all order
func (or *OrderResource) GetAllOrders(ctx *restful.Context) {
	openlog.Info("Got request to get all orders")
	// get query params
	pageno, limit := common.GetPaginationParams(ctx.ReadQueryParameter("pageno"), ctx.ReadQueryParameter("limit"))
	filters := ctx.ReadQueryParameter("filters")
	// call service
	os := order_service.New()
	res := os.GetAllOrders(pageno, limit, filters)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

func (or *OrderResource) OrdersHistory(ctx *restful.Context) {
	openlog.Info("Got request to get all orders history")
	// get query params
	pageno, limit := common.GetPaginationParams(ctx.ReadQueryParameter("pageno"), ctx.ReadQueryParameter("limit"))
	filters := ctx.ReadQueryParameter("filters")
	// call service
	os := order_service.New()
	res := os.OrdersHistory(pageno, limit, filters)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

// get all order
func (or *OrderResource) GetAllDuplicates(ctx *restful.Context) {
	openlog.Info("Got request to get all orders")
	// get query params
	// pageno, limit := common.GetPaginationParams(ctx.ReadQueryParameter("pageno"), ctx.ReadQueryParameter("limit"))
	// filters := ctx.ReadQueryParameter("filters")
	// call service
	os := order_service.New()
	res := os.GetAllDuplicates()
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

// get order by Id
func (or *OrderResource) GetOrderById(ctx *restful.Context) {
	openlog.Info("Got request to get order by id")
	// get id
	id := ctx.ReadPathParameter("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Order ID"}, "application/json")
		return
	}
	// call service
	os := order_service.New()
	res := os.GetOrderById(orderId)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

// get order by Id
func (or *OrderResource) GetOrdersByZone(ctx *restful.Context) {
	openlog.Info("Got request to get order count by zone")
	// get id
	// zoneId := ctx.ReadPathParameter("zoneID")
	// zId, err := strconv.Atoi(zoneId)
	// if err != nil {
	// 	ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Zone ID"}, "application/json")
	// 	return
	// }
	// call service
	os := order_service.New()
	res := os.GetOrdersByZone()
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}
func (or *OrderResource) ConfirmStatus(ctx *restful.Context) {
	// read path params
	openlog.Info("Got a request to change the status by id")
	id := ctx.ReadPathParameter("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Order ID"}, "application/json")
		return
	}
	data := schemas.OrderStatusPayload{}
	ctx.ReadEntity(&data)
	//call service
	os := order_service.New()
	res := os.ConfirmOrderById(orderId, data)
	//return response
	ctx.WriteJSON(res, "application/json")
}

func (or *OrderResource) SelectOrderForOptimization(ctx *restful.Context) {
	openlog.Info("Got request to select order for optimization")
	// get id
	id := ctx.ReadPathParameter("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Order ID"}, "application/json")
		return
	}
	data := schemas.SelectOrder{}
	ctx.ReadEntity(&data)
	//call service
	os := order_service.New()
	res := os.SelectOrderForOptimization(orderId, data)
	//return response
	ctx.WriteJSON(res, "application/json")
}

func (or *OrderResource) UnSelectOrderForOptimization(ctx *restful.Context) {
	openlog.Info("Got request to select order for optimization")
	// get id
	id := ctx.ReadPathParameter("id")
	OrderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Planner ID"}, "application/json")
		return
	}
	//data := schemas.SelectOrder{}
	//ctx.ReadEntity(&data)
	//call service
	os := order_service.New()
	res := os.UnSelectOrderForOptimization(OrderId)
	//return response
	ctx.WriteJSON(res, "application/json")
}

func (or *OrderResource) SelectAllOrdersForOptimization(ctx *restful.Context) {
	openlog.Info("Got request to select all orders for optimization")
	openlog.Info("Got request to get all jobs by plannerdate")
	// get id
	order := schemas.SelectAllOrdersPayload{}
	ctx.ReadEntity(&order)
	// call service
	os := order_service.New()
	res := os.SelectAllOrderForOptimization(order)
	// return response
	ctx.WriteJSON(res, "application/json")
}

// update order
func (or *OrderResource) UpdateOrder(ctx *restful.Context) {
	openlog.Info("Got request to update order")
	// get id
	id := ctx.ReadPathParameter("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Order ID"}, "application/json")
		return
	}
	// read payload
	order := schemas.OrderPayload{}
	ctx.ReadEntity(&order)
	// call service
	os := order_service.New()
	res := os.UpdateOrderById(orderId, order)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

// delete order
func (or *OrderResource) DeleteOrder(ctx *restful.Context) {
	openlog.Info("Got request to delete order")
	// get id
	id := ctx.ReadPathParameter("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.WriteJSON(common.HTTPResponse{Status: 400, Msg: "Invalid Order ID"}, "application/json")
		return
	}
	// call service
	os := order_service.New()
	res := os.DeleteOrderById(orderId)
	// return response
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

func (or *OrderResource) BulkInsert(ctx *restful.Context) {
	openlog.Info("Received request to insert ")
	startTime := time.Now()
	file_data := make([]byte, 0)
	file_name := ""
	config := bulk_upload.BulkInsertConfig{}
	// read mulitpart form data.
	multipart, err := ctx.ReadRequest().MultipartReader()
	if err != nil {
		openlog.Error("Error occured while reading the multipart data.")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
		return
	}
	for {
		part, err := multipart.NextPart()
		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "file":
			file_name = part.FileName()
			file_data, err = ioutil.ReadAll(part)
			fmt.Println("file_data", len(file_data))
			fmt.Println(ctx.Req.Request.Header.Get("Content-Length"))
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "autoInsertNewCustomer":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.AutoInsertCustomer, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "autoInsertNewAddress":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.AutoInsertAddress, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "autoInsertNewItem":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.AutoInsertItems, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "autoDetermineLatLng":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.DeterminzeLatLng, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "autoDetermineZone":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.DertermineZone, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "rejectDuplicate":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.RejectDuplicate, err = strconv.ParseBool(string(ac[:]))
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		}
	}
	if file_name == "" {
		openlog.Error("Error occured while reading the multipart data.")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
		return
	}
	os := order_service.New()
	res := os.BulkInsert(file_data, file_name, config)
	endTime := time.Since(startTime)
	fmt.Printf("Execution time: %v\n", endTime)
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
}

func (or *OrderResource) PullOrdersFromVF(ctx *restful.Context) {
	openlog.Info("Got request to pull orders from VF")

	payload := struct {
		Date   string `json:"date"`
		Status string `json:"status"`
	}{}

	err := ctx.ReadEntity(&payload)
	if err != nil {
		openlog.Error("Error while reading payload " + err.Error())
		res := common.HTTPResponse{Status: 400, Msg: "Error while reading payload"}
		ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
	}

	// Convert struct to JSON string
	jsonString, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// call service
	os := order_service.New()
	res := os.PullOrdersFromVF(jsonString)
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")

}

func (or *OrderResource) PushOrdersToVF(ctx *restful.Context) {
	openlog.Info("Got request to push orders to VF")

	payload := struct {
		PlannerID int `json:"planner_id"`
	}{}

	err := ctx.ReadEntity(&payload)
	if err != nil {
		openlog.Error("Error while reading payload " + err.Error())
		res := common.HTTPResponse{Status: 400, Msg: "Error while reading payload"}
		ctx.WriteHeaderAndJSON(res.Status, res, "application/json")
	}

	// Convert struct to JSON string
	jsonString, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// call service
	os := order_service.New()
	res := os.PushOrdersToVF(jsonString)
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")

}

// URL patterns
func (or *OrderResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: "GET", Path: common.Basepath + "/customer/order", ResourceFunc: or.GetAllOrders, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "GET", Path: common.Basepath + "/customer/order/{id}", ResourceFunc: or.GetOrderById, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "DELETE", Path: common.Basepath + "/customer/order/{id}", ResourceFunc: or.DeleteOrder, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "POST", Path: common.Basepath + "/customer/order", ResourceFunc: or.CreateOrder, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "PUT", Path: common.Basepath + "/customer/order/{id}", ResourceFunc: or.UpdateOrder, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "GET", Path: common.Basepath + "/customer/orders-history", ResourceFunc: or.OrdersHistory, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "POST", Path: common.Basepath + "/orders/bulk-insert", ResourceFunc: or.BulkInsert, Consumes: []string{"multipart/form-data"}, Produces: common.API_HEADERS},

		{Method: "GET", Path: common.Basepath + "/order/zone/count", ResourceFunc: or.GetOrdersByZone, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "POST", Path: common.Basepath + "/order/{id}/_confirm", ResourceFunc: or.ConfirmStatus, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "PUT", Path: common.Basepath + "/orders/{id}/_select", ResourceFunc: or.SelectOrderForOptimization, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "PUT", Path: common.Basepath + "/orders/{id}/_unselect", ResourceFunc: or.UnSelectOrderForOptimization, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "PUT", Path: common.Basepath + "/order/_selectall", ResourceFunc: or.SelectAllOrdersForOptimization, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
		{Method: "GET", Path: common.Basepath + "/orders/_duplicates", ResourceFunc: or.GetAllDuplicates, Consumes: common.API_HEADERS, Produces: common.API_HEADERS},
	}
}
