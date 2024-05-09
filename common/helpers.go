/**
 * Contains reusable common helper functions
 *
**/

package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

// errors
var ErrResourceNotFound = errors.New("resource not found")
var ErrInternalServer = errors.New("internal server error")
var ErrBadRequest = errors.New("bad request")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrConflict = errors.New("conflict")
var ErrJobConflict = errors.New("cannot delete orders with associated jobs")
var ErrDuplicateRecords = errors.New("duplicate records found in database")

const DUPLICATEKEYVALUE = "duplicate"
const EXCELSHEETERROR = "sheet Sheet1 does not exist"

// Order status
const OrderStatusInvalid = "Invalid"
const OrderStatusPending = "Pending"
const OrderStatusConfirmed = "Confirmed"
const OrderStatusScheduled = "Scheduled"
const OrderStatusDelieverd = "Delivered"
const OrderStatusCancelled = "Cancelled"
const OrderStatusReturned = "Returned"
const OrderStatusInTransit = "In Transit"

// Job Status
const JobPending = "Pending"
const JobInProgress = "In Progress"
const JobInTransit = "In Transit"
const JobCompleted = "Completed"
const JobFailed = "Failed"
const JobCancelled = "Cancelled"

// Planner Status

const PlanPending = "Pending"
const PlanOptimized = "Optimized"
const PlanPublished = "Published"
const PlanCancelled = "Cancelled"
const PlanInProgress = "In Progress"
const PlanCompleted = "Completed"

// Job Type
const JobTypePickup = "pickup"
const JobTypeDelivery = "delivery"
const JobTypeStart = "start"
const JobTypeEnd = "end"

// order Type
const OrderDelivery = "Delivery"
const OrderPickup = "Pickup"
const OrderDeliveryAndPickup = "Delivery and Pickup"

// Helper functions

// GetPaginationParams returns the pagination params
func GetPaginationParams(pageno, limit string) (int, int) {
	if pageno == "" {
		pageno = "1" // set to default parameter
	}
	if limit == "" {
		limit = "10" // set to default parameter
	}
	p, err := strconv.Atoi(pageno)
	if err != nil {
		p = 1 // set to default parameter
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		l = 10
	}
	return p, l
}

type PayloadError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func ValidateStruct(payload interface{}) []PayloadError {
	var errors []PayloadError
	// validate
	validate := validator.New()
	err := validate.Struct(payload)
	println(err == nil)
	if err == nil {
		return nil
	}
	for _, err := range err.(validator.ValidationErrors) {
		var payloadError PayloadError
		payloadError.Field = err.Field()
		payloadError.Error = err.Tag()
		errors = append(errors, payloadError)
	}
	return errors
}

func QueryBuilder(addrs []string) string {
	query_start := "SELECT DISTINCT(addr) FROM (VALUES ('"
	query_end := "') ) AS input_addresses(addr) WHERE addr NOT IN (SELECT address_str FROM addresses);"
	query := strings.Join(addrs, "'),('")
	fmt.Println(query)
	return query_start + query + query_end
}

var PlannerSummaryQuery = `
	WITH vehicle_details AS (SELECT driver_vehicle_assigns.id              AS driver_vehicle_assign_id,
		COUNT(DISTINCT jobs.order_id)                                      AS total_orders_per_vehicle,
		SUM(CASE WHEN job_type = 'delivery' THEN orders.weight ELSE 0 END) AS total_weight_per_vehicle
	FROM driver_vehicle_assigns
		LEFT JOIN jobs ON driver_vehicle_assigns.id = jobs.driver_vehicle_assign_id
		LEFT JOIN orders ON jobs.order_id = orders.id
	WHERE driver_vehicle_assigns.planner_id = ?
	GROUP BY driver_vehicle_assigns.id),
	unassigned_orders AS (SELECT sum(weight) AS total_unassigned_weight, COUNT(*) AS total_orders_unassigned
	FROM jobs
			left join orders on jobs.order_id = orders.id
	WHERE jobs.planner_id = ?
	AND job_type = 'delivery'
	AND jobs.driver_vehicle_assign_id IS NULL
	AND jobs.deleted_at IS NULL)
	SELECT COUNT(DISTINCT orders.id)                            AS total_orders_assigned,
	COUNT(DISTINCT jobs.driver_vehicle_assign_id)               AS total_vehicles_assigned,
	(SELECT total_orders_unassigned FROM unassigned_orders)     AS total_orders_unassigned,
	(SELECT total_unassigned_weight FROM unassigned_orders)     AS total_weight_unassigned,
	(SELECT count(*)
	from vehicle_details
	WHERE vehicle_details.total_orders_per_vehicle = 0)         AS total_vehicles_unassigned,
	SUM(estimated_duration)                                     AS total_time,
	SUM(CASE WHEN job_type = 'delivery' THEN weight ELSE 0 END) AS total_weight,
	MAX(weight)                                                 AS max_order_weight,
	MIN(weight)                                                 AS min_order_weight,
	MIN(total_orders_per_vehicle)                               AS min_orders_by_vehicle,
	MAX(total_orders_per_vehicle)                               AS max_orders_by_vehicle,
	MIN(total_weight_per_vehicle)                               AS min_weight_by_vehicle,
	MAX(total_weight_per_vehicle)                               AS max_weight_by_vehicle,
	ROUND(COUNT(DISTINCT orders.id) * 1.0 / COUNT(DISTINCT jobs.driver_vehicle_assign_id),
	2)                                                    AS average_orders_per_vehicle,
	ROUND(SUM(CASE WHEN job_type = 'delivery' THEN weight ELSE 0 END) * 1.0 /
	COUNT(DISTINCT jobs.driver_vehicle_assign_id),
	2)                                                    AS average_weight_per_vehicle

	FROM jobs
	LEFT JOIN orders ON jobs.order_id = orders.id
	LEFT JOIN vehicle_details ON jobs.driver_vehicle_assign_id = vehicle_details.driver_vehicle_assign_id
	WHERE jobs.planner_id = ?
	AND jobs.driver_vehicle_assign_id IS NOT NULL
	AND jobs.deleted_at IS NULL;
	`

var DuplicateOrdersQuery = `select order_no from orders where order_no in ?`
