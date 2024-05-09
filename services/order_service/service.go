package order_service

import (
	"bytes"
	"lp_customer_portal/bulk_upload"
	"lp_customer_portal/common"
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	"lp_customer_portal/repository/customer_repo"
	"lp_customer_portal/repository/item_repo"
	"lp_customer_portal/repository/order_repo"

	//planner_repo "lp_customer_portal/repository/planner_repo"
	//zone_repository "lp_customer_portal/repository/zone_repo"
	"lp_customer_portal/schemas"
	"lp_customer_portal/services/customer_services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
)

type OrderService struct {
	OrderRepo       order_repo.OrderRepoInterface
	CustomerRepo    customer_repo.CustomerRepositoryInterface
	ItemRepo        item_repo.ItemRepositoryInterface
	CustomerService customer_services.CustomerServiceInterface
}

func New() *OrderService {
	openlog.Info("Initializing Order Service")
	db := database.GetClient()
	OrderRepo := order_repo.OrderRepo{DB: db}
	ItemRepo := item_repo.ItemRepo{DB: db}
	CustomerRepo := customer_repo.CustomerRepository{DB: db}

	return &OrderService{
		OrderRepo:    &OrderRepo,
		ItemRepo:     &ItemRepo,
		CustomerRepo: &CustomerRepo,
	}
}

func (service OrderService) CreateOrder(order schemas.OrderPayload) common.HTTPResponse {
	openlog.Debug("Creating the order in the service")
	// check if the order no already exists
	_, err := service.OrderRepo.GetOrderByOrderNo(order.OrderNo)
	if err != nil && err != common.ErrResourceNotFound {
		openlog.Error("Error while fetching the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	// check if customer already exists, if not existed, create a new customer
	if order.CustomerID == 0 { // assumes UI is requesting for creating new customer
		res := service.CustomerService.CreateCustomer(order.Customer) // calling service layer to create customer.
		if res.Status != 200 {
			return common.HTTPResponse{Status: 500, Msg: "Error while creating the customer"}
		}
		order.CustomerID = res.Data.(models.Customer).ID
	}
	// check if address exists, if not create a new address
	if order.AddressID == 0 { // assumes UI is requesting for creating new address
		res := service.CustomerService.CreateAddress(order.Address) // calling service layer to create address.
		if res.Status != 200 {
			return common.HTTPResponse{Status: 500, Msg: "Error while creating the address"}
		}
		order.AddressID = res.Data.(models.Address).ID
	}
	// create the order and order Items
	orderModel := models.Order{}
	order.ToModel(&orderModel)
	// TODO: added zone to order model from address.zone
	// orderModel.ZoneID = order.Address.ZoneID
	// default status
	orderModel.Status = common.OrderStatusPending
	// call repo
	orderModel, err = service.OrderRepo.Insert(orderModel)
	if err != nil {
		openlog.Error("Error while creating the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while creating the order"}
	}
	// return the order no
	return common.HTTPResponse{Status: 200, Msg: "Order created successfully", Data: orderModel}
}

// Get all order
func (service OrderService) GetAllOrders(pageno, limit int, filters string) common.HTTPResponse {
	openlog.Debug("Fetching all orders")
	res, err := service.OrderRepo.GetAllOrders(pageno, limit, filters)
	if err != nil {
		openlog.Error("Error occured while fetching the orders")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching the orders"}
	}
	count, err := service.OrderRepo.GetCountByFilters(filters)
	if err != nil {
		openlog.Error("Error occured while fetching all Orders")
		return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	}
	exists, err := service.OrderRepo.GetISDuplicate()
	if err != nil {
		openlog.Error("Error occured while fetching all Orders")
		return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	}

	IAexists, err := service.OrderRepo.GetInvalidAddress()
	if err != nil {
		openlog.Error("Error occured while fetching all Orders")
		return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	}

	OrderID := []uint{}
	if strings.Contains(filters, "preferred_delivery_start_date") && strings.Contains(filters, "preferred_delivery_end_date") {
		orders, err := service.OrderRepo.GetAllOrders(2, -1, filters)
		if err != nil {
			openlog.Error("Error occured while fetching the orders")
			return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching the orders"}
		}
		for _, order := range orders {
			OrderID = append(OrderID, order.ID)
		}
	}

	data := struct {
		Data           []models.Order `json:"data"`
		Total          int64          `json:"total"`
		IsDuplicate    bool           `json:"is_duplicate"`
		OrderID        []uint         `json:"order_id"`
		InvalidAddress int64          `json:"invalid_address"`
	}{
		Data:           res,
		Total:          count,
		IsDuplicate:    exists,
		OrderID:        OrderID,
		InvalidAddress: IAexists,
	}
	return common.HTTPResponse{Status: 200, Msg: "Orders fetched successfully", Data: data}
}

func (service OrderService) OrdersHistory(pageno, limit int, filters string) common.HTTPResponse {
	openlog.Debug("Fetching all orders")
	res, err := service.OrderRepo.GetAllHistoryOrders(pageno, limit, filters)
	if err != nil {
		openlog.Error("Error occured while fetching the orders")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching the orders"}
	}
	count, err := service.OrderRepo.CountHistoryOrders(filters)
	if err != nil {
		openlog.Error("Error occured while fetching all Orders")
		return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	}
	// exists, err := service.OrderRepo.GetISDuplicate()
	// if err != nil {
	// 	openlog.Error("Error occured while fetching all Orders")
	// 	return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	// }

	// IAexists, err := service.OrderRepo.GetInvalidAddress()
	// if err != nil {
	// 	openlog.Error("Error occured while fetching all Orders")
	// 	return common.HTTPResponse{Msg: "Error occured while fetching Orders", Status: 500}
	// }

	// OrderID := []uint{}
	// if strings.Contains(filters, "preferred_delivery_start_date") && strings.Contains(filters, "preferred_delivery_end_date") {
	// 	orders, err := service.OrderRepo.GetAllOrders(2, -1, filters)
	// 	if err != nil {
	// 		openlog.Error("Error occured while fetching the orders")
	// 		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching the orders"}
	// 	}
	// 	for _, order := range orders {
	// 		OrderID = append(OrderID, order.ID)
	// 	}
	// }

	data := struct {
		Data  []models.Order `json:"data"`
		Total int64          `json:"total"`
	}{
		Data:  res,
		Total: count,
	}

	return common.HTTPResponse{Status: 200, Msg: "Orders fetched successfully", Data: data}
}

func (service OrderService) GetOrderById(id int) common.HTTPResponse {
	openlog.Debug("Fetching the order in the service")
	// call repo
	order, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		openlog.Error("Error while fetching the order" + err.Error())
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Order fetched successfully", Data: order}
}

func (service OrderService) GetOrdersByZone() common.HTTPResponse {
	openlog.Debug("Fetching the order in the service")
	// call repo
	order, err := service.OrderRepo.GetOrderByZoneCount()
	if err != nil {
		openlog.Error("Error while fetching the order" + err.Error())
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order count"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Orders count fetched successfully", Data: order}
}

// update order by id
func (service OrderService) UpdateOrderById(id int, order schemas.OrderPayload) common.HTTPResponse {
	openlog.Debug("Updating the order in the service")
	// check if the order no already exists
	_, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		openlog.Error("Error while fetching the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	// update order
	orderModel := models.Order{}
	err = order.ToModel(&orderModel)
	if err != nil {
		openlog.Error("Error while converting the order payload to model" + err.Error())
		return common.HTTPResponse{Status: 400, Msg: "Error while converting the order payload to model"}
	}
	// call repo
	orderModel, err = service.OrderRepo.UpdateOrderById(id, orderModel)
	if err != nil {
		openlog.Error("Error while updating the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while updating the order"}
	}
	// return the order no
	return common.HTTPResponse{Status: 200, Msg: "Order updated successfully", Data: orderModel}
}

func (service OrderService) ConfirmOrderById(id int, order schemas.OrderStatusPayload) common.HTTPResponse {
	openlog.Debug("conforming the order in the service")
	// check if the order no already exists
	Orderdata, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		openlog.Error("Error while fetching the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	currentTime := time.Now()

	orderModel := models.Order{}

	// order.ToModel(&orderModel)
	if order.PreferredDateTime != "" {
		preffered_date_time, err := time.Parse(common.TIME_FORMAT, order.PreferredDateTime)
		if err != nil {
			openlog.Error("Preferred date should be greater than current date")
			return common.HTTPResponse{Msg: "Preferred date should be greater than current date", Status: 500}
		}
		orderModel.PreferredDeliveryStartDate = preffered_date_time
	}
	// Compare the PrefferedDate with the current time
	if orderModel.PreferredDeliveryStartDate.Sub(currentTime) < 0 {
		openlog.Error("Preferred date should be greater than current date")
		return common.HTTPResponse{Msg: "Preferred date should be greater than current date", Status: 500}

	}

	if order.PreferredEndDateTime != "" {
		preffered_end_date_time, err := time.Parse(common.TIME_FORMAT, order.PreferredEndDateTime)
		if err != nil {
			openlog.Error("Preferred date should be greater than current date")
			return common.HTTPResponse{Msg: "Preferred date should be greater than current date", Status: 500}
		}
		orderModel.PreferredDeliveryEndDate = preffered_end_date_time
	}

	//Compare the PrefferedEndDate with the StartTime
	if orderModel.PreferredDeliveryEndDate.Sub(orderModel.PreferredDeliveryStartDate) < 0 {
		openlog.Error("Preferred end date should be greater than start date")
		return common.HTTPResponse{Msg: "Preferred end date should be greater than start date", Status: 500}
	}

	if Orderdata.Status == common.OrderStatusPending || Orderdata.Status == common.OrderStatusReturned {
		orderModel.Status = common.OrderStatusConfirmed
	}

	// update order
	_, err = service.OrderRepo.UpdateOrderById(id, orderModel)
	if err != nil {
		openlog.Error("Error occurred while conforming order")
		return common.HTTPResponse{Msg: "Error occurred while conforming order", Status: 500}
	}
	// order confirmed successfully
	return common.HTTPResponse{Msg: "order confirmed successfully", Data: orderModel}
}

func (service OrderService) SelectOrderForOptimization(id int, data schemas.SelectOrder) common.HTTPResponse {
	openlog.Debug("Selecting order for optimization")
	// get order by id
	_, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		openlog.Error("Error while fetching the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	//orderModel := models.Order{IsScheduled: data.IsSelected}
	// update order

	order, err := service.OrderRepo.SelectOrderForOptimization(id, data.IsSelected)
	if err != nil {
		openlog.Error("Error updating the order" + err.Error())
		return common.HTTPResponse{Msg: "Error occurred while updating order", Status: 500}
	}
	// order confirmed successfully
	return common.HTTPResponse{Msg: "order Updated successfully", Data: order}
}

//TODO: Unselect order for optimization

func (service OrderService) UnSelectOrderForOptimization(id int) common.HTTPResponse {
	openlog.Debug("UnSelecting order for optimization")
	// get planner by id
	_, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		openlog.Error("Error while fetching the order" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while fetching the order"}
	}
	//orderModel := models.Order{IsScheduled: data.IsSelected}
	// update
	var isSelected *uint

	order, err := service.OrderRepo.UnSelectOrderForOptimization(id, isSelected)
	if err != nil {
		openlog.Error("Error updating the order" + err.Error())
		return common.HTTPResponse{Msg: "Error occurred while updating order", Status: 500}
	}
	// order confirmed successfully
	return common.HTTPResponse{Msg: "order Updated successfully", Data: order}
}

// Delete order by id
func (service OrderService) DeleteOrderById(id int) common.HTTPResponse {
	openlog.Debug("Deleting the order in the service")
	// check if the order no already exists and not deleted
	order, err := service.OrderRepo.GetOrderById(id)
	if err != nil {
		openlog.Error("Error occured while getting order details")
		if err == common.ErrResourceNotFound {
			openlog.Error("Order not found")
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while deleting order"}
	}
	// check if order is already archived.
	if order.Archived {
		openlog.Debug("Order already deleted.")
		return common.HTTPResponse{Status: 409, Msg: "Order already deleted"}
	}
	// update order with is archived true
	update := models.Order{Archived: true}
	order, err = service.OrderRepo.UpdateOrderById(id, update)
	if err != nil {
		openlog.Debug("Error occured while deleting the order")
		return common.HTTPResponse{Status: 500, Msg: "Error Occured while deleting the order"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Deleted order successfully"}
}

// func (service OrderService) DeleteOrderById(id int) common.HTTPResponse {
// 	openlog.Debug("Deleting the order in the service")

// 	// Check if the order exists
// 	order, err := service.OrderRepo.GetOrderById(id)
// 	if err != nil {
// 		openlog.Error("Error occurred while getting order details")
// 		if errors.Is(err, common.ErrResourceNotFound) {
// 			openlog.Error("Order not found")
// 			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
// 		}
// 		openlog.Error("Error occurred while getting order details")
// 		return common.HTTPResponse{Status: 500, Msg: "Error occurred while deleting order"}
// 	}

// 	// Check if a job exists for the order
// 	jobCount, err := service.OrderRepo.GetJobCountByOrderID(id)
// 	if err != nil {
// 		openlog.Error("Error occurred while checking job count")
// 		return common.HTTPResponse{Status: 500, Msg: "Error occurred while deleting order"}
// 	}

// 	if jobCount != nil && *jobCount > 0 {
// 		// Update order to set Archived to true
// 		order.Archived = true
// 		_, err = service.OrderRepo.UpdateOrderById(id, order)
// 		if err != nil {
// 			openlog.Debug("Error occurred while updating the order")
// 			return common.HTTPResponse{Status: 500, Msg: "Error occurred while deleting order"}
// 		}
// 	} else {
// 		// Delete order by ID
// 		_, err := service.OrderRepo.DeleteOrderById(id)
// 		if err != nil {
// 			openlog.Debug("Error occurred while deleting the order")
// 			return common.HTTPResponse{Status: 500, Msg: "Error occurred while deleting order"}
// 		}
// 	}

// 	return common.HTTPResponse{Status: 200, Msg: "Deleted order successfully"}
// }

func (service OrderService) BulkInsert(fileData []byte, filename string, config bulk_upload.BulkInsertConfig) common.HTTPResponse {
	openlog.Debug("inserting bulk orders")
	inserter := bulk_upload.NewBulkInsert(config)

	// Determine file type
	fileType := bulk_upload.GetFileType(filename)

	// Create a reader based on file type
	var records [][]string
	var err error
	switch fileType {
	case "csv":
		records, err = bulk_upload.ReadCSV(fileData)
		if err != nil {
			openlog.Error("Error while reading the file: " + err.Error())
			return common.HTTPResponse{Status: 400, Msg: "Error while reading the file"}
		}
	case "xlsx":
		records, err = bulk_upload.ReadExcel(fileData)
		if err != nil {
			openlog.Error("Error while reading the file: " + err.Error())
			if err.Error() == common.EXCELSHEETERROR {
				return common.HTTPResponse{Status: 400, Msg: "Error while reading excel file: Sheet1 not found"}
			}
			return common.HTTPResponse{Status: 400, Msg: "Error while reading the file"}
		}
	default:
		openlog.Error("Unsupported file type for " + filename)
		return common.HTTPResponse{Status: 400, Msg: "Unsupported file type for " + filename}
	}

	count, err := inserter.ProcessOrderRecords(records)
	if err != nil {
		return common.HTTPResponse{Status: 500, Msg: err.Error()}
	}
	return common.HTTPResponse{
		Status: 200,
		Msg:    "Inserted orders of count " + strconv.FormatInt(count, 10) + "",
	}
}

func (service OrderService) SelectAllOrderForOptimization(data schemas.SelectAllOrdersPayload) common.HTTPResponse {
	openlog.Debug("Selecting order for optimization")

	orderModel := models.Order{PlannerID: data.IsSelected}
	// update order
	// fmt.Println(preffered_start_date, preffered_end_date)
	_, err := service.OrderRepo.UpdateOrderForOtimizationById(data.PreferredStartDate, data.PreferredEndDate, orderModel)
	if err != nil {
		openlog.Error("Error updating the order")
		if err.Error() == "no orders found for the date" {
			return common.HTTPResponse{Msg: "No orders found for the date", Status: 404}
		}
		return common.HTTPResponse{Msg: "Error occurred while updating order", Status: 500}
	}
	// order confirmed successfully
	return common.HTTPResponse{Msg: "order Updated successfully", Status: 200}
}

// Get all order
func (service OrderService) GetAllDuplicates() common.HTTPResponse {
	openlog.Debug("Fetching all orders")
	duplicates, err := service.OrderRepo.GetAllDuplicates()
	if err != nil {
		openlog.Error("Error occured while fetching the orders")
		return common.HTTPResponse{Status: 500, Msg: "Error occured while fetching the orders"}
	}
	duplicateOrderNos := make([]string, 0)
	orderNos := make([]string, 0)
	for _, o_no := range duplicates {
		duplicateOrderNos = append(duplicateOrderNos, o_no.OrderNo)
	}

	for _, dup_order := range duplicateOrderNos {
		o_no := strings.Split(dup_order, "_")[0]
		orderNos = append(orderNos, o_no)
	}
	duplicateOrderNos = append(duplicateOrderNos, orderNos...)
	return common.HTTPResponse{Status: 200, Msg: "Orders fetched successfully", Data: duplicateOrderNos}
}

func (service OrderService) PullOrdersFromVF(payload []byte) common.HTTPResponse {
	openlog.Debug("Fetching orders from VF")
	// http request to fastapi
	VF_URL := archaius.GetString("versafleet.pullOrders", "")

	req, err := http.NewRequest("POST", VF_URL, bytes.NewReader(payload))
	if err != nil {
		openlog.Error("Error while creating request: " + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while creating request"}
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		openlog.Error("Error while pulling VF data" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while pulling VF data"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Orders pulled successfully from VF"}
}

func (service OrderService) PushOrdersToVF(payload []byte) common.HTTPResponse {
	openlog.Debug("Pushing orders to VF")
	// http request to fastapi
	VF_URL := archaius.GetString("versafleet.pushOrders", "")

	req, err := http.NewRequest("POST", VF_URL, bytes.NewReader(payload))
	if err != nil {
		openlog.Error("Error while creating request: " + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while creating request"}
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		openlog.Error("Error while pushing data to VF" + err.Error())
		return common.HTTPResponse{Status: 500, Msg: "Error while pushing data to VF"}
	}
	return common.HTTPResponse{Status: 200, Msg: "Orders pushed successfully to VF"}
}

func (service OrderService) DeleteOrdersByIds(id []int) common.HTTPResponse {
	openlog.Debug("Deleting all orders")
	_, err := service.OrderRepo.DeleteOrdersByOrderIds(id)
	if err != nil {
		openlog.Error("Error occured while deleting the orders")
		if err == common.ErrJobConflict {
			return common.HTTPResponse{Status: 400, Msg: "cannot delete orders with associated jobs"}
		}
		if err == common.ErrResourceNotFound {
			openlog.Error("Order not found")
			return common.HTTPResponse{Status: 404, Msg: "Order not found"}
		}
		return common.HTTPResponse{Status: 500, Msg: "Error occured while deleting order"}
	}

	return common.HTTPResponse{Status: 200, Msg: "Deleted order successfully"}

}
