/**
 * Service Layer implementation goes here.
 *
**/

package item_services

import (
	//"lp_oms/models"

	"lp_customer_portal/bulk_upload"
	common "lp_customer_portal/common"
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	item_repo "lp_customer_portal/repository/item_repo"
	"lp_customer_portal/schemas"
	"strconv"

	"github.com/go-chassis/openlog"
)

type ItemService struct {
	ItemRepo item_repo.ItemRepositoryInterface
}

func New() *ItemService {
	openlog.Info("Initializing Item Service")
	db := database.GetClient()
	ItemRepo := item_repo.ItemRepo{DB: db}
	return &ItemService{
		ItemRepo: &ItemRepo,
	}
}

// FetchAll service implementation
func (rs *ItemService) FetchAllItems(pageno, limit int, filters string) common.HTTPResponse {
	res, err := rs.ItemRepo.FetchAllItems(pageno, limit, filters) // call repository layer
	if err != nil {
		openlog.Error("Error occured while fetching all items")
		return common.HTTPResponse{Msg: "Error occured while fetching items", Status: 500}
	}
	count, err := rs.ItemRepo.FetchAllItemsCount(filters)
	if err != nil {
		openlog.Error("Error occurred while fetching all items")
		return common.HTTPResponse{Msg: "Error occurred while fetching items", Status: 500}
	}
	data := struct {
		Data  []models.Item `json:"data"`
		Total int64         `json:"total"`
	}{
		Data:  res, // Use the converted items slice here
		Total: count,
	}
	return common.HTTPResponse{Msg: "Items fetched successfully", Data: data, Status: 200}
}

// FetchUserById service implementation
func (rs *ItemService) FetchItemById(id int) common.HTTPResponse {
	openlog.Debug("Fetching item by id ")
	res, err := rs.ItemRepo.FetchItemById(id) // call repository layer
	if err != nil {
		openlog.Error("Error occured while fetching all Item")
		return common.HTTPResponse{Msg: "Error occured while fetching Item" + err.Error()}
	}
	return common.HTTPResponse{Msg: "Item fetched successfully", Data: res, Status: 202}
}

func (rs *ItemService) BulkInsertItems(fileData []byte, filename string, config bulk_upload.BulkInsertConfig) common.HTTPResponse {
	openlog.Debug("inserting bulk items")
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

	count, err := inserter.ProcessItemRecords(records)
	if err != nil {
		return common.HTTPResponse{Status: 500}
	}
	return common.HTTPResponse{
		Status: 200,
		Msg:    "Inserted items of count " + strconv.FormatInt(count, 10) + "",
	}
}
func (rs *ItemService) CreateItem(item schemas.ItemPayload) common.HTTPResponse {
	openlog.Debug("Got a request to create item")

	// Check if an item with the same Name or SKU already exists
	existingItem, err := rs.ItemRepo.GetItemByNameOrSKU(item.Name, item.SKU)
	if err != nil {
		openlog.Error("Error occurred while checking existing items")
		return common.HTTPResponse{Msg: "Error occurred while creating item", Status: 500}
	}
	// If an existing item is found, return HTTP 409 Conflict status
	if existingItem != (models.Item{}) {
		return common.HTTPResponse{Msg: "An item with the same name or SKU already exists", Status: 409}
	}
	// convert from schema to model
	modelItem := models.Item{}
	item.ToModel(&modelItem)

	data, err := rs.ItemRepo.Insert(modelItem) // call repository layer
	if err != nil {
		openlog.Error("Error occurred while creating item")
		return common.HTTPResponse{Msg: "Error occurred while creating item", Status: 500}
	}
	return common.HTTPResponse{Msg: "item created successfully", Data: data, Status: 201}
}

func (rs *ItemService) UpdateItemById(id int, item schemas.ItemPayload) common.HTTPResponse {
	openlog.Debug("Updating item by id " + string(rune(id)))
	// check if item exists
	res, err := rs.ItemRepo.GetById(id)
	if err != nil {
		openlog.Error("Error occurred while fetching item")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Msg: "item not found", Status: 404}
		}
		return common.HTTPResponse{Msg: "Error occurred while updating item", Status: 500}
	}
	// update the item.
	res.Name = item.Name
	res.Description = item.Description
	res.Category = item.Category
	res.SKU = item.SKU
	res.Capacity = item.Capacity
	res.Volume = item.Volume
	res.Weight = item.Weight
	res.InnerPack = item.InnerPack
	res.Status = item.Status

	data, err := rs.ItemRepo.UpdateById(id, res) // call repository layer
	if err != nil {
		openlog.Error("Error occurred while updating item")
		return common.HTTPResponse{Msg: "Error occurred while updating item", Status: 500}
	}
	// item updated successfully
	return common.HTTPResponse{Msg: "item updated successfully", Data: data, Status: 200}
}

func (rs *ItemService) DeleteItemById(id int) common.HTTPResponse {
	// check if item exists
	_, err := rs.ItemRepo.GetById(id)
	if err != nil {
		openlog.Error("Error occurred while fetching item")
		if err == common.ErrResourceNotFound {
			return common.HTTPResponse{Msg: "item not found", Status: 404}
		}
		return common.HTTPResponse{Msg: "Error occurred while deleting item", Status: 500}
	}
	// delete the item.
	err = rs.ItemRepo.DeleteById(id) // call repository layer
	if err != nil {
		openlog.Error("Error occurred while deleting item")
		return common.HTTPResponse{Msg: "Error occurred while deleting item", Status: 500}
	}
	// item deleted successfully
	return common.HTTPResponse{Msg: "item deleted successfully", Status: 200}
}
