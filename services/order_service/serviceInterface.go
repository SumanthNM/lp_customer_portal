package order_service

import (
	"lp_customer_portal/bulk_upload"
	"lp_customer_portal/common"
	"lp_customer_portal/schemas"
)

type OrderServiceInterface interface {
	CreateOrder(order schemas.OrderPayload) common.HTTPResponse
	GetOrderById(orderId int) common.HTTPResponse
	GetAllOrders(pageno, limit int, filters string) common.HTTPResponse
	OrdersHistory(pageno, limit int, filters string) common.HTTPResponse
	UpdateOrderById(orderId int, order schemas.OrderPayload) common.HTTPResponse
	DeleteOrderById(orderId int) common.HTTPResponse
	BulkInsert(fileData []byte, filename string, config bulk_upload.BulkInsertConfig) common.HTTPResponse
	ConfirmOrderById(id int, order schemas.OrderStatusPayload) common.HTTPResponse
	SelectOrderForOptimization(id int, data schemas.SelectOrder) common.HTTPResponse
	SelectAllOrderForOptimization(data schemas.SelectAllOrdersPayload) common.HTTPResponse
	UnSelectOrderForOptimization(id int) common.HTTPResponse
	GetOrdersByZone() common.HTTPResponse
	GetAllDuplicates() common.HTTPResponse
	PullOrdersFromVF(payload []byte) common.HTTPResponse
	PushOrdersToVF(payload []byte) common.HTTPResponse
	DeleteOrdersByIds(id []int) common.HTTPResponse
}
