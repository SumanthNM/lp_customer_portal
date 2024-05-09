package order_repo

import (
	"errors"
	"lp_customer_portal/common"
	"lp_customer_portal/models"
	"strings"
	"time"

	"github.com/go-chassis/openlog"
	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

func (or OrderRepo) Insert(order models.Order) (models.Order, error) {
	openlog.Debug("Inserting order into database with order ID " + order.OrderNo)
	res := or.DB.Create(&order)
	return order, res.Error
}

// Get All Orders
func (or OrderRepo) GetAllOrders(pageno, limit int, filters string) ([]models.Order, error) {
	openlog.Debug("Fetching all orders")
	filterScope := common.GetAllCondition(filters)
	orders := make([]models.Order, 0)
	res := or.DB.Scopes(filterScope...).Model(&models.Order{}).Order("preferred_pickup_start_date DESC").Offset((pageno - 1) * limit).Limit(limit).Joins("DeliveryAddress").Joins("PickupAddress").Joins("Customer").Preload("OrderItems").Where("archived != true").Find(&orders)
	return orders, res.Error
}

func (or OrderRepo) GetAllHistoryOrders(pageno, limit int, filters string) ([]models.Order, error) {
	openlog.Debug("Fetching all orders")
	filterScope := common.GetAllCondition(filters)
	// Add condition for PreferredDeliveryEndDate
	currentDate := time.Now().UTC()
	filterScope = append(filterScope, func(db *gorm.DB) *gorm.DB {
		return db.Where("preferred_delivery_end_date >= ?", currentDate)
	})
	orders := make([]models.Order, 0)
	res := or.DB.Scopes(filterScope...).Model(&models.Order{}).Order("preferred_pickup_start_date DESC").Offset((pageno - 1) * limit).Limit(limit).Joins("DeliveryAddress").Joins("PickupAddress").Joins("Customer").Preload("OrderItems").Where("archived != true").Find(&orders)
	return orders, res.Error
}

func (or OrderRepo) CountHistoryOrders(filters string) (int64, error) {
	// Construct the filter scope
	filterScope := common.GetAllCondition(filters)

	// Add condition for PreferredDeliveryEndDate
	currentDate := time.Now().UTC()
	filterScope = append(filterScope, func(db *gorm.DB) *gorm.DB {
		return db.Where("preferred_delivery_end_date >= ?", currentDate)
	})

	// Count orders with the provided filters
	var count int64
	res := or.DB.Scopes(filterScope...).
		Model(&models.Order{}).
		Where("archived != true").
		Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

// Get All Duplicates
func (or OrderRepo) GetAllDuplicates() ([]models.Order, error) {
	openlog.Debug("Fetching all orders")
	orders := make([]models.Order, 0)
	res := or.DB.Model(&models.Order{}).Where("duplicates = ?", true).Order("order_no").Find(&orders)
	return orders, res.Error
}

func (or OrderRepo) GetInvalidAddress() (int64, error) {
	openlog.Info("Fetching count of Invalid address from database")
	var count int64
	res := or.DB.Model(&models.Order{}).Where("status = ?", common.OrderStatusInvalid).Where("deleted_at IS NULL").Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	// Return the count of orders with an "Invalid" status
	return count, nil
}

// Get Order BY Id
func (or OrderRepo) GetOrderById(id int) (models.Order, error) {
	openlog.Debug("Fetching order by id ")
	order := models.Order{}
	res := or.DB.Model(&models.Order{}).Preload("DeliveryAddress").Preload("PickupAddress").Preload("Customer").Preload("OrderItems").Preload("OrderItems.Item").Where("id = ? AND archived != true", id).First(&order)
	if res.Error != nil {
		openlog.Error("Error while fetching order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

// Update Order by Id
func (or OrderRepo) UpdateOrderById(id int, order models.Order) (models.Order, error) {
	openlog.Debug("Updating order by id ")
	res := or.DB.Model(&models.Order{}).Where("id = ?", id).Updates(order) // will falses
	if res.Error != nil {
		openlog.Error("Error while updating order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

func (or OrderRepo) UpdateOrderForOtimizationById(preferred_start_date, preferred_end_date string, order models.Order) ([]models.Order, error) {
	openlog.Debug("Updating order by id ")
	orders := make([]models.Order, 0)
	//fmt.Println(preferred_start_date, preferred_end_date)
	res := or.DB.Model(&models.Order{}).Where("status = ?", "Confirmed").Where("preferred_start_date >= ?", preferred_start_date).Where("preferred_start_date <= ?", preferred_end_date).Find(&orders)
	if res.Error != nil {
		openlog.Error("Error while updating order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
		return orders, res.Error
	}
	// fmt.Println(orders)
	if (len(orders)) == 0 {
		return orders, errors.New("no orders found for the date")
	}
	updateRes := or.DB.Model(&models.Order{}).Where("status = ?", "Confirmed").Where("preferred_start_date >= ?", preferred_start_date).Where("preferred_start_date <= ?", preferred_end_date).Update("planner_id", order.IsSelected)
	if updateRes.Error != nil {
		openlog.Error("Error while updating order by id " + updateRes.Error.Error())
		if updateRes.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
	}
	return orders, nil
}

// Delete Order By Id
func (or OrderRepo) DeleteOrderById(id int) (models.Order, error) {
	openlog.Debug("Deleting duplicate orders")
	order := models.Order{}
	res := or.DB.Model(&models.Order{}).Where("id = ?", id).Delete(&order)
	if res.Error != nil {
		openlog.Error("Error while deleting order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

// Delete Order By Id
func (or OrderRepo) DeleteDuplicateOrders() ([]models.Order, error) {
	openlog.Debug("Deleting order by id ")
	orders := make([]models.Order, 0)
	res := or.DB.Unscoped().Model(&models.Order{}).Where("duplicates = ?", true).Delete(&orders)
	if res.Error != nil {
		openlog.Error("Error while deleting order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
		return orders, res.Error
	}
	return orders, nil
}

// Delete Orders By Id
func (or OrderRepo) DeleteOrdersByIds(order []string) ([]models.Order, error) {
	openlog.Debug("Deleting order by id ")
	orders := make([]models.Order, 0)
	res := or.DB.Unscoped().Model(&models.Order{}).Where("order_no in ?", order).Delete(&orders)
	if res.Error != nil {
		openlog.Error("Error while deleting order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
		return orders, res.Error
	}
	return orders, nil
}

// Get Order by order no
func (or OrderRepo) GetOrderByOrderNo(orderNo string) (models.Order, error) {
	openlog.Debug("Fetching order by order no ")
	order := models.Order{}
	res := or.DB.Model(&models.Order{}).Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		openlog.Error("Error while fetching order by order no " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

func (or OrderRepo) BulkInsert(orders []models.Order) (int64, error) {
	openlog.Debug("Bulk inserting orders")
	// fmt.Println(orders)
	// chuncks
	chunkSize := 100
	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		res := or.DB.Create(orders[i:end])
		if res.Error != nil {
			openlog.Error("Error occured while inserting the orders")
			if strings.Contains(res.Error.Error(), common.DUPLICATEKEYVALUE) {
				return 0, common.ErrDuplicateRecords
			}
			return 0, res.Error
		}
	}

	return int64(len(orders)), nil
}

func (or OrderRepo) GetCount() (int64, error) {
	openlog.Info("Fetching count of orders from database")
	var count int64
	result := or.DB.Model(&models.Order{}).Where("archived != true").Count(&count)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, result.Error
	}
	return count, nil
}

func (pr OrderRepo) GetCountByFilters(filters string) (int64, error) {
	openlog.Info("Fetching count of orders from database")
	var count int64
	filterScope := common.GetAllCondition(filters)
	result := pr.DB.Scopes(filterScope...).Model(&models.Order{}).Joins("DeliveryAddress").Joins("PickupAddress").Joins("Customer").Count(&count)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, result.Error
	}
	return count, nil
}

func (or OrderRepo) GetCountbyDate(plannerID int, fromdate, todate, status string, pageno, limit int, filters string) (int64, int64, error) {
	openlog.Info("Fetching count of orders from database")
	var count int64
	result := or.DB.Model(&models.Order{}).Where("preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ? AND status = ?", fromdate, todate, status).Or("preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ? AND status = ?", fromdate, todate, status).Count(&count)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, 0, result.Error
	}
	var selectedCount int64
	result = or.DB.Model(&models.Order{}).Where("preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ? AND status = ?", fromdate, todate, status).Or("preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ? AND status = ?", fromdate, todate, status).Count(&selectedCount)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, 0, result.Error
	}
	return count, selectedCount, nil
}

func (or OrderRepo) GetOrdersByDateStatus(fromdate, todate, status string, pageno, limit int, filters string) ([]models.Order, error) {
	openlog.Info("Fetching orders from database by date and status")
	orders := make([]models.Order, 0)
	result := or.DB.Model(&models.Order{}).Offset((pageno - 1) * limit).Limit(limit).Preload("DeliveryAddress").Preload("DeliveryAddress.Zone").Preload("Customer")
	result = result.Where("preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ? AND status = ?", fromdate, todate, status).Or("preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ? AND status = ?", fromdate, todate, status).Order("id DESC").Preload("DeliveryAddress").Preload("PickupAddress").Find(&orders)
	if result.Error != nil {
		openlog.Error("Error occured while fetching orders from database by date and status")
		return orders, result.Error
	}
	return orders, nil
}

func (or OrderRepo) GetSelectedConfirmedOrders(plannerID int) ([]models.Order, error) {
	openlog.Info("Fetching selected orders from database")
	ordersModel := make([]models.Order, 0)
	result := or.DB.Model(&models.Order{}).Preload("Address").Preload("Address.Zone").Preload("Customer").Preload("OrderItems").Where("planner_id = ? ", plannerID).Find(&ordersModel)
	if result.Error != nil {
		openlog.Error("Error occured while fetching selected orders from database")
		return ordersModel, result.Error
	}
	return ordersModel, nil
}

func (or OrderRepo) BulkUpdateOrders(orders []models.Order) (int64, error) {
	openlog.Debug("Bulk updating orders")
	// chuncks
	chunkSize := 100
	var rowAffected int64 = 0
	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		res := or.DB.Save(orders[i:end])
		if res.Error != nil {
			openlog.Error("Error occured while updating the orders")
			return 0, res.Error
		}
		rowAffected += res.RowsAffected
	}
	return rowAffected, nil
}

func (or OrderRepo) SelectOrderForOptimization(order int, selected uint) (models.Order, error) {
	openlog.Debug("Updating order by id ")
	orderModel := models.Order{}
	res := or.DB.Model(orderModel).Where("id = ?", order).Update("planner_id", selected).Find(&orderModel)
	if res.Error != nil {
		openlog.Error("Error while updating order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orderModel, common.ErrResourceNotFound
		}
		return orderModel, res.Error
	}
	return orderModel, nil
}

func (or OrderRepo) UnSelectOrderForOptimization(order int, selected *uint) (models.Order, error) {
	openlog.Debug("Updating order by id ")
	orderModel := models.Order{}
	res := or.DB.Model(orderModel).Where("id = ?", order).Update("planner_id", selected).Find(&orderModel)
	if res.Error != nil {
		openlog.Error("Error while updating order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orderModel, common.ErrResourceNotFound
		}
		return orderModel, res.Error
	}
	return orderModel, nil
}

func (or OrderRepo) UpdatePlannerIDByOrderDate(plannerID int, status string, startDate time.Time, endDate time.Time) (int64, error) {
	openlog.Debug("Updating order by id ")
	var orders []models.Order
	updateRes := or.DB.Model(&models.Order{}).Where("status = ?", status).Where("preferred_start_date >= ?", startDate).Where("preferred_start_date <= ?", endDate).Update("planner_id", plannerID)
	if updateRes.Error != nil {
		openlog.Error("Error while updating order by id " + updateRes.Error.Error())
		if updateRes.Error == gorm.ErrRecordNotFound {
			return 0, common.ErrResourceNotFound
		}
	}
	return int64(len(orders)), nil
}

func (or OrderRepo) GetOrderByZoneCount() (int64, error) {
	openlog.Info("Fetching orders from the database")
	var count int64
	query := `
	SELECT COUNT(o.order_no) AS order_count
    FROM orders o
    LEFT JOIN addresses a ON o.address_id = a.id
    GROUP BY COALESCE(a.zone_id, 0);
	`
	result := or.DB.Raw(query).Scan(&count)
	if result.Error != nil {
		openlog.Error("Error occurred while fetching orders from the database")
		return 0, result.Error // You should return 0 in case of an error.
	}

	return count, nil
}

func (or OrderRepo) GetOrdersByIds(ids []int64) ([]models.Order, error) {
	openlog.Info("Fetching orders from the database")
	orders := make([]models.Order, 0)
	result := or.DB.Model(&models.Order{}).Preload("PickupAddress").Preload("PickupAddress.Zone").Preload("DeliveryAddress").Preload("DeliveryAddress.Zone").Preload("Customer").Preload("OrderItems").Where("id IN ?", ids).Find(&orders)
	if result.Error != nil {
		openlog.Error("Error occurred while fetching orders from the database")
		return orders, result.Error // You should return 0 in case of an error.
	}

	return orders, nil
}

func (or OrderRepo) PublishOrders(plannerId int) ([]models.Order, error) {
	openlog.Debug("Unpublish orders by list of id ")
	order := []models.Order{}
	orderMap := map[string]interface {
	}{
		"PlannerID": plannerId,
		"Status":    common.OrderStatusScheduled,
	}
	subquery := or.DB.Select("order_id").Where("planner_id = ?", plannerId).Where("deleted_at is null").Table("jobs")
	res := or.DB.Model(&models.Order{}).Where("id in (?)", subquery).Updates(orderMap)
	if res.Error != nil {
		openlog.Error("Error while fetching order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

func (or OrderRepo) UnpublishOrders(plannerId int) ([]models.Order, error) {
	openlog.Debug("Unpublish orders by list of id ")
	order := []models.Order{}
	orderMap := map[string]interface {
	}{
		"PlannerID": nil,
		"Status":    common.OrderStatusConfirmed,
	}

	res := or.DB.Model(&models.Order{}).Where("planner_id = ?", plannerId).Updates(orderMap)
	if res.Error != nil {
		openlog.Error("Error while fetching order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return order, common.ErrResourceNotFound
		}
		return order, res.Error
	}
	return order, nil
}

func (or OrderRepo) GetUnselectOrdersByDateStatus(fromdate, todate, status string, plannerid int, pageno, limit int, filters string) ([]models.Order, error) {
	openlog.Info("Fetching orders from database by date and status")
	orders := make([]models.Order, 0)
	query := `
	SELECT * from orders
    WHERE ((preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ?) OR (preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ?)) AND status = ?
    AND id NOT IN (select order_id from jobs
    WHERE planner_id = ? and job_type = 'delivery' and jobs.deleted_at is null)
	order by id desc
	`
	if limit > 0 && pageno > 0 {
		query = query + "limit ? offset ?"
		result := or.DB.Raw(query, fromdate, todate, fromdate, todate, status, plannerid, limit, (pageno-1)*limit)
		result = result.Preload("DeliveryAddress").Preload("DeliveryAddress.Zone")
		result = result.Preload("PickupAddress").Preload("PickupAddress.Zone")
		result = result.Preload("Customer")

		result = result.Find(&orders)
		if result.Error != nil {
			openlog.Error("Error occured while fetching orders from database by date and status")
			return orders, result.Error
		}
		return orders, nil
	} else {
		result := or.DB.Raw(query, fromdate, todate, fromdate, todate, status, plannerid)
		result = result.Preload("DeliveryAddress").Preload("DeliveryAddress.Zone")
		result = result.Preload("PickupAddress").Preload("PickupAddress.Zone")
		result = result.Preload("Customer")
		result = result.Find(&orders)
		if result.Error != nil {
			openlog.Error("Error occured while fetching orders from database by date and status")
			return orders, result.Error
		}
		return orders, nil
	}
}

func (or OrderRepo) GetUnselectCountbyDate(plannerID int, fromdate, todate, status string, pageno, limit int, filters string) (int64, int64, error) {
	openlog.Info("Fetching count of orders from database")

	var count int64
	query := `
	SELECT COUNT(*) from orders
    WHERE ((preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ?) OR (preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ?)) AND status = ?
    AND id NOT IN (select order_id from jobs
    WHERE planner_id = ? and job_type = 'delivery' and jobs.deleted_at is null
    )`
	result := or.DB.Raw(query, fromdate, todate, fromdate, todate, status, plannerID).Count(&count)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, 0, result.Error
	}
	var selectedCount int64
	result = or.DB.Model(&models.Order{}).Where("preferred_pickup_start_date >= ? AND preferred_pickup_end_date <= ? AND status = ? AND planner_id = ?", fromdate, todate, status, plannerID).Or("preferred_delivery_start_date >= ? AND preferred_delivery_end_date <= ? AND status = ? AND planner_id = ?", fromdate, todate, status, plannerID).Count(&selectedCount)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of orders from database")
		return count, 0, result.Error
	}
	return count, selectedCount, nil
}

func (or OrderRepo) FetchDuplicatedOrderNo(query string, orderNo []string) ([]string, error) {
	var orderNoList []string

	res := or.DB.Raw(query, orderNo).Find(&orderNoList)
	if res.Error != nil {
		openlog.Error("Error while fetching duplicated orders " + res.Error.Error())
		return []string{}, res.Error
	}
	return orderNoList, nil
}

// Delete Orders By Id
func (or OrderRepo) DeleteOrdersByOrderIds(order []int) ([]models.Order, error) {
	openlog.Debug("Deleting order by id ")
	// Check if order_ids exist in jobs table
	var jobsCount int64
	if err := or.DB.Model(&models.Jobs{}).Where("order_id IN ?", order).Count(&jobsCount).Error; err != nil {
		openlog.Error("Error while checking if order_ids exist in jobs table: " + err.Error())
		return nil, err
	}

	// If order_ids exist in jobs table, return an error or handle accordingly
	if jobsCount > 0 {
		return nil, common.ErrJobConflict
	}

	orders := make([]models.Order, 0)
	res := or.DB.Unscoped().Model(&models.OrderItem{}).Where("order_id in ?", order).Delete(&models.OrderItem{})
	if res.Error != nil {
		openlog.Error("Error while deleting order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
		return orders, res.Error
	}
	res = or.DB.Unscoped().Model(&models.Order{}).Where("id in ?", order).Delete(&orders)
	if res.Error != nil {
		openlog.Error("Error while deleting order by id " + res.Error.Error())
		if res.Error == gorm.ErrRecordNotFound {
			return orders, common.ErrResourceNotFound
		}
		return orders, res.Error
	}
	return orders, nil
}

func (or OrderRepo) GetISDuplicate() (bool, error) {
	openlog.Info("Fetching count of isduplicates from database")
	var count int64
	res := or.DB.Model(&models.Order{}).Where("duplicates = ?", true).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	// If count is greater than 1, it means there are duplicates
	return count > 0, nil
}

// GetJobCountByOrderID returns the count of jobs associated with a specific order ID
func (or OrderRepo) GetJobCountByOrderID(orderID int) (*int64, error) {
	var count int64
	result := or.DB.Model(&models.Jobs{}).Where("order_id = ?", orderID).Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	return &count, nil
}
