// /**
//  * Insert Bulk order based on the configuration.
//  *
// **/

package bulk_upload

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"lp_oms/common"
// 	"lp_oms/database"
// 	"lp_oms/here_helper"
// 	"lp_oms/latlng_helper"
// 	"lp_oms/models"
// 	optimizer_common "lp_oms/optimizer/common"
// 	"math"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/go-chassis/openlog"
// 	"golang.org/x/time/rate"
// )

// /*
// Bulk Insertion from CSV
// 1. Preload all the following information from DB to prevent multiple db calls
// 	-Items
// 	-Customers
// 	-Zones
// 	-Addresses
// 2. Iterate through CSV to obtain all the rows of data
// 3. Extract all the new, distinct addresses from the data by comparing with the db
// 	-Call AWS Location Service API to fetch latlng
// 	-Update zone with the latlng
// 	-bulk insert new, distinct addresses into addresses table
// 4. Bulk upload all the items
// 5. Bulk upload all the customers
// 6. Preload all the data from the db w the newly inserted info
// 7. Iterate through all the rows of data
// 	-combined orders with the same orderNo  [given data has one item per row. one order can have multiple items ]
// 	-update all the columns in the Order
// 8. Bulk upload the orders

// NOTE:
// 	-start & end time are hardcoded to 09:00 - 17:00
// 	-weight & volumne of items are hardcoded to 10
// */

// func (bi *BulkInsert) preloadMetaData() error {
// 	// preload items
// 	items, err := bi.ItemRepo.FetchAllItems(2, -1, "")
// 	if err != nil {
// 		openlog.Error("Error while preloading items")
// 		return errors.New("error while preloading items")
// 	}
// 	bi.Items = items
// 	// preload customers
// 	res, err := bi.CustomerRepo.GetAllCustomers(2, -1, "")
// 	if err != nil {
// 		openlog.Error("Error while preloading customers")
// 		return errors.New("error while preloading customers")
// 	}
// 	bi.Customers = res

// 	// preload zones
// 	zones, err := bi.ZoneRepo.GetAll(2, -1, true)
// 	if err != nil {
// 		openlog.Error("Error while preloading zones")
// 		return errors.New("error while preloading zones")
// 	}
// 	bi.Zones = zones

// 	// preload addresses
// 	addresses, err := bi.CustomerRepo.GetAllAddress(2, -1)
// 	if err != nil {
// 		openlog.Error("Error while preloading zones")
// 		return errors.New("error while preloading zones")
// 	}
// 	bi.Addresses = addresses
// 	return nil
// }

// // will insert new customer if not found here.
// func (bi *BulkInsert) findCustomer(order OrderRow) (models.Customer, error) {
// 	for _, customer := range bi.Customers {
// 		if customer.CustomerNo == order.CustomerNo {
// 			return customer, nil
// 		}
// 	}
// 	return models.Customer{}, errors.New("customer not found")
// }

// // find the item id with sku
// func (bi *BulkInsert) findItem(orderRow OrderRow) (models.Item, error) {
// 	for _, item := range bi.Items {
// 		if item.SKU == orderRow.ItemRow.ItemSKU {
// 			return item, nil
// 		}
// 	}
// 	return models.Item{}, errors.New("item not found")
// }

// // find address by address name
// func (bi *BulkInsert) findAddressByName(orderRow OrderRow) (models.Address, error) {
// 	for _, address := range bi.Addresses {
// 		if address.Address_Str == orderRow.AddressStr {
// 			return address, nil
// 		}
// 	}
// 	return models.Address{}, errors.New("address not found")
// }

// // LoadCSV reads the csv, users the key map to do data transformation into OrderRow
// func (b *BulkInsert) ProcessOrderRecords(records [][]string) (int64, error) {
// 	openlog.Debug("Loading the csv ")
// 	openlog.Debug("Starting the transaction")

// 	database.StartTransaction()
// 	err := b.preloadMetaData()
// 	if err != nil {
// 		openlog.Error("Error while preloading the metadata")
// 		return 0, errors.New("error occured while preloading the metadata")
// 	}

// 	// get headers
// 	headers := records[0]
// 	// remove header from remaining data
// 	records = records[1:]

// 	orderRows := make([]OrderRow, 0)

// 	for _, record := range records {
// 		data := make(map[string]string)
// 		for i := 0; i < len(record); i++ {
// 			data[headers[i]] = record[i]
// 		}
// 		order, err := transform(data)
// 		if err != nil {
// 			openlog.Error("Error while transforming the data")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, err
// 		}
// 		orderRows = append(orderRows, *order)
// 	}
// 	openlog.Debug("Successfully read the csv ")
// 	openlog.Debug("Number of orders read [" + strconv.Itoa(len(orderRows)) + "]")
// 	b.OrderRows = orderRows

// 	// 1. fetch all unique address, items and customers from csv. check w db to see if there are any new ones.
// 	var addressList = []string{}
// 	var itemModelsList = []models.Item{}
// 	var customerModelList = []models.Customer{}
// 	var orderNoList = []string{}
// 	uniqueItem := make(map[string]bool)
// 	uniqueCustomer := make(map[string]bool)
// 	// need to do a deduplication
// 	for _, OrderRow := range b.OrderRows {
// 		addressList = append(addressList, OrderRow.AddressStr)
// 		orderNoList = append(orderNoList, OrderRow.OrderNo)
// 		itemSKU := OrderRow.ItemRow.ItemSKU
// 		if _, exists := uniqueItem[itemSKU]; !exists {
// 			uniqueItem[itemSKU] = true
// 			item := getItem(OrderRow.ItemRow)
// 			itemModelsList = append(itemModelsList, item)
// 		}
// 		customer := OrderRow.CustomerNo
// 		if _, exists := uniqueCustomer[customer]; !exists {
// 			uniqueCustomer[customer] = true
// 			customer := getCustomer(OrderRow)
// 			customerModelList = append(customerModelList, customer)
// 		}
// 	}
// 	query := common.QueryBuilder(addressList)
// 	uniqueAddresses, err := b.CustomerRepo.FetchUniqueAddress(query)
// 	if err != nil {
// 		openlog.Error("Error while fetching all new unique addresses")
// 		openlog.Debug("Rolling back the transaction")
// 		database.RollbackTransaction()
// 		return 0, err
// 	}

// 	duplicateOrderNo, err := b.OrderRepo.FetchDuplicatedOrderNo(common.DuplicateOrdersQuery, orderNoList)
// 	if err != nil {
// 		openlog.Error("Error while fetching duplicated orderNo")
// 		openlog.Debug("Rolling back the transaction")
// 		database.RollbackTransaction()
// 		return 0, err
// 	}
// 	// check if there are new addresses to be inserted into the db
// 	if len(uniqueAddresses) != 0 {
// 		// 2. fetch all address info. latlng & zone
// 		var addressModels []models.Address
// 		if b.Config.DeterminzeLatLng {
// 			addressModels = b.FetchLatLngFromAddr(uniqueAddresses)
// 		}
// 		// 3. upsert into address table
// 		if b.Config.AutoInsertAddress {
// 			count, err := b.CustomerRepo.BulkInsertAddresses(addressModels)
// 			if err != nil {
// 				openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
// 				openlog.Debug("Rolling back the transaction")
// 				database.RollbackTransaction()
// 				return 0, err
// 			}
// 			openlog.Debug("Number of addresses inserted [" + strconv.FormatInt(count, 10) + "]")
// 		}
// 	}
// 	// bulk upload new items
// 	if b.Config.AutoInsertItems {
// 		count, err := b.ItemRepo.BulkInsertItems(itemModelsList)
// 		if err != nil {
// 			openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, err
// 		}
// 		openlog.Debug("Number of items inserted [" + strconv.FormatInt(count, 10) + "]")
// 	}
// 	// bulk upload customers
// 	if b.Config.AutoInsertCustomer {
// 		count, err := b.CustomerRepo.BulkInsertCustomer(customerModelList)
// 		if err != nil {
// 			openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, err
// 		}
// 		openlog.Debug("Number of customers inserted [" + strconv.FormatInt(count, 10) + "]")
// 	}
// 	err = b.preloadMetaData()
// 	if err != nil {
// 		openlog.Error("Error while preloading the metadata")
// 		return 0, errors.New("error occured while preloading the metadata")
// 	}

// 	// convert order rows into orders
// 	orders := make(map[string]*models.Order)

// 	// Rule is one row - on order
// 	for _, orderRow := range orderRows {
// 		order := &models.Order{}

// 		// check if the order already exists in map
// 		if _, ok := orders[orderRow.OrderNo]; ok {
// 			order = orders[orderRow.OrderNo]
// 		}

// 		order.OrderNo = orderRow.OrderNo
// 		// order.Priority = orderRow.Priority
// 		order.OrderDate = orderRow.OrderedDate
// 		customer, err := b.findCustomer(orderRow)
// 		if err != nil {
// 			openlog.Error("Error while finding the customer")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, errors.New("error occured while finding the customer")
// 		}

// 		order.CustomerID = customer.ID
// 		order.Customer = customer

// 		// search for address by unique address name
// 		address, err := b.findAddressByName(orderRow)
// 		if err != nil {
// 			openlog.Error("Error while finding the address")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, errors.New("error occured while finding the address")
// 		}
// 		order.AddressID = address.ID
// 		order.ZoneID = address.ZoneID
// 		order.Address = &address

// 		order.PreferredStartDate = orderRow.PreferredDeliverStartTime
// 		order.PreferredEndDate = orderRow.PreferredDeliverEndTime
// 		order.Status = common.OrderStatusConfirmed // default status
// 		// find item in db
// 		item, err := b.findItem(orderRow)
// 		if err != nil {
// 			openlog.Error("Item not found")
// 			openlog.Debug("Rolling back the transaction")
// 			database.RollbackTransaction()
// 			return 0, err
// 		}
// 		orderItemWeight := 0.0
// 		orderItemVolume := 0.0
// 		// if QuantityMain = 0, take QuantityMinor/Innerpack * weight
// 		if orderRow.QuantityMain == 0 {
// 			orderItemWeight = math.Ceil(item.Weight * float64(orderRow.QuantityMinor) / float64(item.InnerPack))
// 			orderItemVolume = math.Ceil(item.Volume * float64(orderRow.QuantityMinor) / float64(item.InnerPack))
// 		} else {
// 			orderItemWeight = math.Ceil(item.Weight * float64(orderRow.QuantityMain))
// 			orderItemVolume = math.Ceil(item.Volume * float64(orderRow.QuantityMain))
// 		}

// 		order.OrderItems = append(order.OrderItems, models.OrderItem{
// 			ItemID:        item.ID,
// 			QuantityMain:  orderRow.QuantityMain,
// 			QuantityMinor: orderRow.QuantityMinor,
// 			Weight:        orderItemWeight,
// 			Volume:        orderItemVolume,
// 		})
// 		order.Weight = order.Weight + orderItemWeight
// 		order.Volume = order.Volume + orderItemVolume
// 		// update volume and weight of same orders w diff items & quantity type
// 		orders[orderRow.OrderNo] = order // adding order to map
// 	}

// 	if len(orders) == 0 {
// 		openlog.Error("No orders found in the csv")
// 		openlog.Debug("Rolling back the transaction")
// 		database.RollbackTransaction()
// 		return 0, errors.New("no orders found in the csv")
// 	}
// 	// convert map to array
// 	orderData := make([]models.Order, 0)
// 	for _, v := range orders {
// 		if v.Address.Latitude < 0 || v.Address.Longitude < 0 {
// 			v.Status = common.OrderStatusInvalid
// 		}
// 		isDuplicate := false
// 		for _, duplicate := range duplicateOrderNo {
// 			if v.OrderNo == duplicate {
// 				isDuplicate = true
// 				v.OrderNo += "_duplicate"
// 				v.Duplicates = true
// 			}
// 		}
// 		if isDuplicate {
// 			if !b.Config.RejectDuplicate {
// 				orderData = append(orderData, *v)
// 			}
// 		} else {
// 			orderData = append(orderData, *v)
// 		}
// 	}
// 	// Insert into database
// 	count, err := b.OrderRepo.BulkInsert(orderData)
// 	if err != nil {
// 		openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
// 		openlog.Debug("rollback the transaction")
// 		database.RollbackTransaction() // rollback the transaction
// 		return 0, err
// 	}
// 	openlog.Debug("commiting the trasaction")
// 	database.CommitTransaction()
// 	openlog.Debug("Number of orders inserted [" + strconv.FormatInt(count, 10) + "]")
// 	return count, nil
// }

// func transform(data map[string]string) (*OrderRow, error) {
// 	order := OrderRow{}
// 	keymap := GetKeyMap("order", "orderKeyMap", ORDERKEYS)
// 	order.OrderNo = data[keymap["orderNo"]]
// 	order.CustomerName = data[keymap["customerName"]]
// 	order.CustomerNo = data[keymap["customerName"]] + data[keymap["channel"]] + data[keymap["customerType"]] // NOTE: using CUSTOMER_NAME + ADDRESS to get unique costomerNo
// 	order.ItemRow.Name = data[keymap["itemName"]]
// 	order.ItemRow.ItemSKU = data[keymap["itemSKU"]]
// 	// hardcoded as a default for missing data points
// 	order.ItemRow.Weight = 10
// 	order.ItemRow.Volume = 10

// 	order.ZipCode = data[keymap["zip"]]
// 	order.AddressStr = data[keymap["deliveryStreetAddress"]]
// 	qtyMain, err := strconv.Atoi(data[keymap["quantityMain"]])
// 	if err != nil {
// 		fmt.Println("error", err.Error())
// 		openlog.Error("Error while converting the quantity main to int")
// 		return nil, errors.New("error occured while parsing quantity main")
// 	}
// 	order.QuantityMain = qtyMain
// 	qtyMinor, err := strconv.Atoi(data[keymap["quantityMinor"]])
// 	if err != nil {
// 		openlog.Error("Error while converting the quantity minor to int")
// 		return nil, errors.New("error occured while parsing quantity minor")
// 	}
// 	order.QuantityMinor = qtyMinor
// 	date, err := parseAndFormatDate(data[keymap["orderedDate"]], EXCEL_DATE_FORMATS)
// 	if err != nil {
// 		openlog.Error("Error while converting string to date [" + err.Error() + "]")
// 		return nil, errors.New("error occured while parsing date")
// 	}
// 	order.OrderedDate = date

// 	date, err = parseAndFormatDate(data[keymap["deliveryDate"]], EXCEL_DATE_FORMATS)
// 	if err != nil {
// 		openlog.Error("Error while converting string to date [" + err.Error() + "]")
// 		return nil, errors.New("error occured while parsing date")
// 	}
// 	// startTime is hardcoded 09:00
// 	preferredStartDateTime := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
// 	order.PreferredDeliverStartTime = preferredStartDateTime
// 	// endTime is hardcoded to 17:00
// 	preferredEndDateTime := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location())
// 	order.PreferredDeliverEndTime = preferredEndDateTime

// 	order.Priority = data[keymap["priority"]]

// 	return &order, nil
// }

// func getCustomer(row OrderRow) models.Customer {
// 	customerObj := models.Customer{
// 		CustomerNo:   row.CustomerNo,
// 		CustomerName: row.CustomerName,
// 	}
// 	return customerObj
// }

// var maxWorkers = 20
// var rateLimit = 200 // aws rate limit

// type Job struct {
// 	Address  string
// 	Location optimizer_common.Location
// 	Error    error
// }

// // this is used to control the rate limitting of aws location service
// var limiter = rate.NewLimiter(rate.Limit(rateLimit), 1)

// func (bi *BulkInsert) FetchLatLngFromAddr(address []models.Address) []models.Address {
// 	var (
// 		jobs           = make(chan models.Address, rateLimit) // jobs queue
// 		resultsChannel = make(chan models.Address, len(address))
// 		wg             sync.WaitGroup
// 	)
// 	// initiation the workers.
// 	for i := 0; i < maxWorkers; i++ { // config or for hardcoded.
// 		wg.Add(1)
// 		go bi.worker(i, jobs, resultsChannel, &wg) // these are threads
// 	}
// 	// Generate HTTP POST requests with addresses as input
// 	for i := 0; i < len(address); i++ {
// 		job := address[i]
// 		if err := limiter.Wait(context.Background()); err != nil {
// 			openlog.Error("Rate limit exceeded for worker" + string(rune(i)))
// 			continue
// 		}
// 		jobs <- job
// 	}
// 	close(jobs)
// 	wg.Wait()
// 	close(resultsChannel)
// 	results := make([]models.Address, 0)
// 	for addr := range resultsChannel {
// 		results = append(results, addr)
// 	}
// 	return results
// }

// func (bi *BulkInsert) worker(id int, jobs chan models.Address, resultsChannel chan models.Address, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for job := range jobs {
// 		// Wait for permission from the rate limiter.
// 		lat, lng, err := latlng_helper.GetLocationsFromText(job.Address_Str)
// 		if err != nil {
// 			openlog.Error("Error while converting address to latlng with AWS: " + err.Error())
// 		}
// 		job.Latitude = lat
// 		job.Longitude = lng

// 		if lat == -1 {
// 			lat, lng, err := here_helper.GetCoordinates(job.Address_Str)
// 			job.Latitude = lat
// 			job.Longitude = lng
// 			if err != nil {
// 				openlog.Error("Error while converting address to latlng with HERE: " + err.Error())
// 			}
// 		}

// 		if bi.Config.DertermineZone {
// 			zone, err := common.GetZone(job.Latitude, job.Longitude, bi.Zones)
// 			if err != nil {
// 				openlog.Error("Error while determining the zone")
// 			}
// 			if err == nil {
// 				job.ZoneID = zone.ID
// 			}
// 		}
// 		resultsChannel <- job
// 	}
// }
