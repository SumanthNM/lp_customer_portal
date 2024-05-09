/**
 * Insert Bulk order based on the configuration.
 *
**/

package bulk_upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"lp_customer_portal/common"
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	"math"
	"strconv"

	"github.com/go-chassis/openlog"
)

/*
Bulk Insertion from CSV
1. Preload all the following information from DB to prevent multiple db calls
	-Items
	-Customers
	-Zones
	-Addresses
2. Iterate through CSV to obtain all the rows of data
3. Extract all the new, distinct addresses from the data by comparing with the db
	-Call AWS Location Service API to fetch latlng
	-Update zone with the latlng
	-bulk insert new, distinct addresses into addresses table
4. Bulk upload all the items
5. Bulk upload all the customers
6. Preload all the data from the db w the newly inserted info
7. Iterate through all the rows of data
	-combined orders with the same orderNo  [given data has one item per row. one order can have multiple items ]
	-update all the columns in the Order
8. Bulk upload the orders

NOTE:
	-start & end time are hardcoded to 09:00 - 17:00
	-weight & volumne of items are hardcoded to 10
*/

func (bi *BulkInsert) preloadMetaData() error {
	// preload items
	items, err := bi.ItemRepo.FetchAllItems(2, -1, "")
	if err != nil {
		openlog.Error("Error while preloading items")
		return errors.New("error while preloading items")
	}
	bi.Items = items
	// preload customers
	res, err := bi.CustomerRepo.GetAllCustomers(2, -1, "")
	if err != nil {
		openlog.Error("Error while preloading customers")
		return errors.New("error while preloading customers")
	}
	bi.Customers = res

	// preload zones
	// zones, err := bi.ZoneRepo.GetAll(2, -1, true)
	// if err != nil {
	// 	openlog.Error("Error while preloading zones")
	// 	return errors.New("error while preloading zones")
	// }
	// bi.Zones = zones

	// preload addresses
	addresses, err := bi.CustomerRepo.GetAllAddress(2, -1)
	if err != nil {
		openlog.Error("Error while preloading zones")
		return errors.New("error while preloading zones")
	}
	bi.Addresses = addresses
	return nil
}

// will insert new customer if not found here.
func (bi *BulkInsert) findCustomer(order OrderRow) (models.Customer, error) {
	for _, customer := range bi.Customers {
		if customer.CustomerNo == order.CustomerNo {
			return customer, nil
		}
	}
	return models.Customer{}, errors.New("customer not found")
}

// find the item id with sku
func (bi *BulkInsert) findItem(orderRow OrderRow) (models.Item, error) {
	for _, item := range bi.Items {
		if item.SKU == orderRow.ItemRow.ItemSKU {
			return item, nil
		}
	}
	return models.Item{}, errors.New("item not found")
}

// find address by address name
func (bi *BulkInsert) findAddressByName(addressStr string) (models.Address, error) {
	for _, address := range bi.Addresses {
		if address.Address_Str == addressStr {
			return address, nil
		}
	}
	return models.Address{}, errors.New("address not found")
}

// LoadCSV reads the csv, users the key map to do data transformation into OrderRow
func (b *BulkInsert) ProcessOrderRecords(records [][]string) (int64, error) {
	openlog.Debug("Loading the csv ")
	openlog.Debug("Starting the transaction")

	database.StartTransaction()
	err := b.preloadMetaData()
	if err != nil {
		openlog.Error("Error while preloading the metadata")
		return 0, errors.New("error occured while preloading the metadata")
	}

	// get headers
	headers := records[0]
	// remove header from remaining data
	records = records[1:]

	orderRows := make([]OrderRow, 0)

	for _, record := range records {
		data := make(map[string]string)
		for i := 0; i < len(record); i++ {
			data[headers[i]] = record[i]
		}
		order, err := transform(data)
		if err != nil {
			openlog.Error("Error while transforming the data")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, err
		}
		orderRows = append(orderRows, *order)
	}
	openlog.Debug("Successfully read the csv ")
	openlog.Debug("Number of orders read [" + strconv.Itoa(len(orderRows)) + "]")
	b.OrderRows = orderRows

	// 1. fetch all unique address, items and customers from csv. check w db to see if there are any new ones.
	// var orderNoList = []string{}
	// var itemModelsList = []models.Item{}
	var customerModelList = []models.Customer{}
	// uniqueItem := make(map[string]bool)
	uniqueCustomer := make(map[string]bool)

	uniqueAddresses := make(map[string]bool)
	var addressModelList = []models.Address{}
	// need to do a deduplication
	for _, OrderRow := range b.OrderRows {
		if _, exists := uniqueAddresses[OrderRow.PickupAddress]; !exists {
			uniqueAddresses[OrderRow.PickupAddress] = true
			addressModelList = append(addressModelList, models.Address{
				Address_Str: OrderRow.PickupAddress,
				Latitude:    OrderRow.PickupLat,
				Longitude:   OrderRow.PickupLng,
				ZoneID:      uint(OrderRow.Zone),
				ZipCode:     uint(OrderRow.PickupPostal),
			})
		}
		if _, exists := uniqueAddresses[OrderRow.DeliveryAddress]; !exists {
			uniqueAddresses[OrderRow.DeliveryAddress] = true
			addressModelList = append(addressModelList, models.Address{
				Address_Str: OrderRow.DeliveryAddress,
				Latitude:    OrderRow.DeliveryLat,
				Longitude:   OrderRow.DeliveryLng,
				ZoneID:      uint(OrderRow.Zone),
				ZipCode:     uint(OrderRow.DeliveryPostal),
			})
		}

		// itemSKU := OrderRow.ItemRow.ItemSKU
		// if _, exists := uniqueItem[itemSKU]; !exists {
		// 	uniqueItem[itemSKU] = true
		// 	item := getItem(OrderRow.ItemRow)
		// 	itemModelsList = append(itemModelsList, item)
		// }
		customer := OrderRow.CustomerNo
		if _, exists := uniqueCustomer[customer]; !exists {
			uniqueCustomer[customer] = true
			customer := getCustomer(OrderRow)
			customerModelList = append(customerModelList, customer)
		}
	}

	// duplicateOrderNo, err := b.OrderRepo.FetchDuplicatedOrderNo(common.DuplicateOrdersQuery, orderNoList)
	// if err != nil {
	// 	openlog.Error("Error while fetching duplicated orderNo")
	// 	openlog.Debug("Rolling back the transaction")
	// 	database.RollbackTransaction()
	// 	return 0, err
	// }
	// check if there are new addresses to be inserted into the db
	// if len(uniqueAddresses) != 0 {
	// 2. fetch all address info. latlng & zone
	// var addressModels []models.Address
	// if b.Config.DeterminzeLatLng {
	// 	addressModels = b.FetchLatLngFromAddr(uniqueAddresses)
	// }
	// 	// 3. upsert into address table
	if b.Config.AutoInsertAddress {
		count, err := b.CustomerRepo.BulkInsertAddresses(addressModelList)
		if err != nil {
			openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, err
		}
		openlog.Debug("Number of addresses inserted [" + strconv.FormatInt(count, 10) + "]")
	}
	// }
	// bulk upload new items
	// if b.Config.AutoInsertItems {
	// 	count, err := b.ItemRepo.BulkInsertItems(itemModelsList)
	// 	if err != nil {
	// 		openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
	// 		openlog.Debug("Rolling back the transaction")
	// 		database.RollbackTransaction()
	// 		return 0, err
	// 	}
	// 	openlog.Debug("Number of items inserted [" + strconv.FormatInt(count, 10) + "]")
	// }
	// bulk upload customers
	if b.Config.AutoInsertCustomer {
		count, err := b.CustomerRepo.BulkInsertCustomer(customerModelList)
		if err != nil {
			openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, err
		}
		openlog.Debug("Number of customers inserted [" + strconv.FormatInt(count, 10) + "]")
	}
	err = b.preloadMetaData()
	if err != nil {
		openlog.Error("Error while preloading the metadata")
		openlog.Debug("Rolling back the transaction")
		database.RollbackTransaction()
		return 0, errors.New("error occured while preloading the metadata")
	}

	// convert order rows into orders
	orders := make(map[string]*models.Order)

	// Rule is one row - on order
	for _, orderRow := range orderRows {
		order := &models.Order{}

		// check if the order already exists in map
		if _, ok := orders[orderRow.OrderNo]; ok {
			order = orders[orderRow.OrderNo]
		}

		order.OrderNo = orderRow.OrderNo
		// order.Priority = orderRow.Priority
		order.OrderDate = orderRow.OrderedDate
		customer, err := b.findCustomer(orderRow)
		if err != nil {
			openlog.Error("Error while finding the customer")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, errors.New("error occured while finding the customer")
		}

		order.CustomerID = customer.ID
		order.Customer = customer

		order.OrderType = orderRow.JobType

		// search for address by unique address name
		pickupAddress, err := b.findAddressByName(orderRow.PickupAddress)
		if err != nil {
			openlog.Error("Error while finding the address")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, errors.New("error occured while finding the address")
		}
		order.PickupAddressID = pickupAddress.ID
		order.PickupAddress = &pickupAddress

		deliveryAddress, err := b.findAddressByName(orderRow.DeliveryAddress)
		if err != nil {
			openlog.Error("Error while finding the address")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, errors.New("error occured while finding the address")
		}
		order.DeliveryAddressID = deliveryAddress.ID
		order.DeliveryAddress = &deliveryAddress

		if orderRow.JobType == common.OrderDelivery {
			order.ZoneID = deliveryAddress.ZoneID
		}
		if orderRow.JobType == common.OrderPickup {
			order.ZoneID = pickupAddress.ZoneID
		}

		order.PreferredPickupStartDate = orderRow.PreferredPickupStartDateTime
		order.PreferredPickupEndDate = orderRow.PreferredPickupEndDateTime
		order.PickupServiceTime = orderRow.PickupCustomServiceTime

		order.PreferredDeliveryStartDate = orderRow.PreferredDeliveryStartDateTime
		order.PreferredDeliveryEndDate = orderRow.PreferredDeliveryEndDateTime
		order.DeliveryServiceTime = orderRow.DeliveryCustomServiceTime

		order.SkillsetID = uint(orderRow.Skills)

		order.Status = common.OrderStatusConfirmed // default status
		// find item in db
		item, err := b.findItem(orderRow)
		if err != nil {
			openlog.Error("Item not found")
			openlog.Debug("Rolling back the transaction")
			database.RollbackTransaction()
			return 0, err
		}

		orderItemWeight := math.Ceil(item.Weight * float64(orderRow.QuantityMain))
		orderItemVolume := math.Ceil(item.Volume * float64(orderRow.QuantityMain))

		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ItemID:       item.ID,
			QuantityMain: orderRow.QuantityMain,
			// QuantityMinor: orderRow.QuantityMinor,
			Weight: orderItemWeight,
			Volume: orderItemVolume,
		})
		order.Weight = order.Weight + orderItemWeight
		order.Volume = order.Volume + orderItemVolume
		// update volume and weight of same orders w diff items & quantity type
		orders[orderRow.OrderNo] = order // adding order to map
	}

	if len(orders) == 0 {
		openlog.Error("No orders found in the csv")
		openlog.Debug("Rolling back the transaction")
		database.RollbackTransaction()
		return 0, errors.New("no orders found in the csv")
	}
	// convert map to array
	orderData := make([]models.Order, 0)
	for _, v := range orders {
		if v.DeliveryAddress.Latitude < 0 || v.DeliveryAddress.Longitude < 0 {
			v.Status = common.OrderStatusInvalid
		}
		if v.PickupAddress.Latitude < 0 || v.PickupAddress.Longitude < 0 {
			v.Status = common.OrderStatusInvalid
		}
		orderData = append(orderData, *v)
	}
	// Insert into database
	count, err := b.OrderRepo.BulkInsert(orderData)
	if err != nil {
		openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
		openlog.Debug("rollback the transaction")
		database.RollbackTransaction() // rollback the transaction
		return 0, err
	}
	openlog.Debug("commiting the trasaction")
	database.CommitTransaction()
	openlog.Debug("Number of orders inserted [" + strconv.FormatInt(count, 10) + "]")
	return count, nil
}

func transform(data map[string]string) (*OrderRow, error) {
	order := OrderRow{}
	keymap := GetKeyMap("order", "A21", A21ORDERKEYS)
	order.OrderNo = data[keymap["trackingNumber"]]
	order.CustomerNo = data[keymap["customerName"]]
	order.CustomerName = data[keymap["customerName"]]

	if data[keymap["quantity"]] == "" {
		order.QuantityMain = 1
	} else {
		qtyMain, err := strconv.Atoi(data[keymap["quantity"]])
		if err != nil {
			fmt.Println("error", err.Error())
			openlog.Error("Error while converting the quantity main to int")
			return nil, errors.New("error occured while parsing quantity main")
		}
		order.QuantityMain = qtyMain
	}

	// hard code item into the orders
	order.ItemRow.ItemSKU = "000001"

	date, err := parseAndFormatDate(data[keymap["orderedDate"]], EXCEL_DATE_FORMATS)
	if err != nil {
		openlog.Error("Error while converting string to date [" + err.Error() + "]")
		return nil, errors.New("error occured while parsing date")
	}
	order.OrderedDate = date
	order.PickupAddress = data[keymap["pickupStreetAddress"]]
	pickupPostal, err := strconv.Atoi(data[keymap["pickupPostal"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the pickup postal to int")
		return nil, errors.New("error occured while parsing pickup postal")
	}

	var pickupCoordinates []float64
	err = json.Unmarshal([]byte(data[keymap["pickupLatLng"]]), &pickupCoordinates)
	if err != nil {
		fmt.Println("Error while unmarshaling latlng data", err)
		return nil, err
	}
	order.PickupLat = pickupCoordinates[0]
	order.PickupLng = pickupCoordinates[1]

	order.PickupPostal = pickupPostal

	pickupStartTime, err := parseAndFormatDate(data[keymap["preferredPickupStartDateTime"]], EXCEL_DATE_FORMATS)
	if err != nil {
		openlog.Error("Error while converting string to date [" + err.Error() + "]")
		return nil, errors.New("error occured while parsing date")
	}
	order.PreferredPickupStartDateTime = pickupStartTime
	pickupEndTime, err := parseAndFormatDate(data[keymap["preferredPickupEndDateTime"]], EXCEL_DATE_FORMATS)
	if err != nil {
		openlog.Error("Error while converting string to date [" + err.Error() + "]")
		return nil, errors.New("error occured while parsing date")
	}
	order.PreferredPickupEndDateTime = pickupEndTime
	pickupServiceTime, err := strconv.Atoi(data[keymap["pickupCustomServiceTime"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the service time to int")
		return nil, errors.New("error occured while parsing service time")
	}

	order.PickupCustomServiceTime = pickupServiceTime

	order.DeliveryAddress = data[keymap["deliveryStreetAddress"]]
	deliveryPostal, err := strconv.Atoi(data[keymap["deliveryPostal"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the delivery postal to int")
		return nil, errors.New("error occured while parsing delivery postal")
	}
	order.DeliveryPostal = deliveryPostal

	var deliveryCoordinates []float64
	err = json.Unmarshal([]byte(data[keymap["deliveryLatLng"]]), &deliveryCoordinates)
	if err != nil {
		fmt.Println("Error while unmarshaling latlng data", err)
		return nil, err
	}
	order.DeliveryLat = deliveryCoordinates[0]
	order.DeliveryLng = deliveryCoordinates[1]

	deliveryStartTime, err := parseAndFormatDate(data[keymap["preferredDeliveryStartDateTime"]], EXCEL_DATE_FORMATS)
	if err != nil {
		openlog.Error("Error while converting string to date [" + err.Error() + "]")
		return nil, errors.New("error occured while parsing date")
	}
	order.PreferredDeliveryStartDateTime = deliveryStartTime
	deliveryEndTime, err := parseAndFormatDate(data[keymap["preferredDeliveryEndDateTime"]], EXCEL_DATE_FORMATS)
	if err != nil {
		openlog.Error("Error while converting string to date [" + err.Error() + "]")
		return nil, errors.New("error occured while parsing date")
	}
	order.PreferredDeliveryEndDateTime = deliveryEndTime
	deliveryServiceTime, err := strconv.Atoi(data[keymap["deliveryCustomServiceTime"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the service time to int")
		return nil, errors.New("error occured while parsing service time")
	}

	order.DeliveryCustomServiceTime = deliveryServiceTime

	zone, err := strconv.Atoi(data[keymap["zone"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the zone to int")
		return nil, errors.New("error occured while parsing zone")
	}
	order.Zone = zone

	skills, err := strconv.Atoi(data[keymap["skills"]])
	if err != nil {
		fmt.Println("error", err.Error())
		openlog.Error("Error while converting the skills to int")
		return nil, errors.New("error occured while parsing skills")
	}
	order.Skills = skills

	if data[keymap["jobType"]] == "pickup" {
		order.JobType = common.OrderPickup
	}
	if data[keymap["jobType"]] == "delivery" {
		order.JobType = common.OrderDelivery
	}

	return &order, nil
}

func getCustomer(row OrderRow) models.Customer {
	customerObj := models.Customer{
		CustomerNo:   row.CustomerNo,
		CustomerName: row.CustomerName,
	}
	return customerObj
}
