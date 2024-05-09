package order_repo

import (
	"lp_customer_portal/models"
	"time"
)

type OrderRepoInterface interface {
	Insert(order models.Order) (models.Order, error)
	GetAllOrders(pageno, limit int, filters string) ([]models.Order, error)
	GetAllHistoryOrders(pageno, limit int, filters string) ([]models.Order, error)
	CountHistoryOrders(filters string) (int64, error)
	GetOrderById(id int) (models.Order, error)
	UpdateOrderById(id int, order models.Order) (models.Order, error)
	DeleteOrderById(id int) (models.Order, error)
	GetOrderByOrderNo(orderNo string) (models.Order, error)
	BulkInsert(orders []models.Order) (int64, error) //TODO:
	GetOrdersByDateStatus(fromdate, todate, status string, pageno, limit int, filters string) ([]models.Order, error)
	GetCount() (int64, error)
	GetCountByFilters(filters string) (int64, error)
	GetSelectedConfirmedOrders(plannerID int) ([]models.Order, error)
	GetCountbyDate(plannerID int, fromdate, todate, status string, pageno, limit int, filters string) (int64, int64, error)
	BulkUpdateOrders(orders []models.Order) (int64, error)
	UpdateOrderForOtimizationById(preferred_start_date, preferred_end_date string, order models.Order) ([]models.Order, error)
	SelectOrderForOptimization(order int, selected uint) (models.Order, error)
	UnSelectOrderForOptimization(order int, selected *uint) (models.Order, error)
	UpdatePlannerIDByOrderDate(plannerID int, status string, startDate time.Time, endDate time.Time) (int64, error) // updates planner_id in orders table for orders between start and end date
	GetOrderByZoneCount() (int64, error)
	GetOrdersByIds(ids []int64) ([]models.Order, error)
	PublishOrders(plannerId int) ([]models.Order, error)
	UnpublishOrders(plannerId int) ([]models.Order, error)
	GetUnselectOrdersByDateStatus(fromdate, todate, status string, plannerid int, pageno, limit int, filters string) ([]models.Order, error)
	GetUnselectCountbyDate(plannerID int, fromdate, todate, status string, pageno, limit int, filters string) (int64, int64, error)
	FetchDuplicatedOrderNo(query string, orderNo []string) ([]string, error)
	GetAllDuplicates() ([]models.Order, error)
	DeleteDuplicateOrders() ([]models.Order, error)
	DeleteOrdersByIds(order []string) ([]models.Order, error)
	DeleteOrdersByOrderIds(order []int) ([]models.Order, error)
	GetISDuplicate() (bool, error)
	GetInvalidAddress() (int64, error)
	GetJobCountByOrderID(orderID int) (*int64, error)
}
