package item_repo

import (
	"errors"
	"lp_customer_portal/common"
	"lp_customer_portal/models"

	"github.com/go-chassis/openlog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemRepo struct {
	DB *gorm.DB
}

func (ir ItemRepo) FetchAllItems(pageno, limit int, filters string) ([]models.Item, error) {
	openlog.Debug("fetching all the items.")
	items := []models.Item{}
	filterScope := common.GetAllCondition(filters)
	err := ir.DB.Scopes(filterScope...).Model(&models.Item{}).Offset((pageno - 1) * limit).Limit(limit).Find(&items).Error // fmt.Println(items)
	if err != nil {
		openlog.Error("error occured while fetching all the items." + err.Error())
		return items, errors.New("error occured while fetching all the items")
	}
	return items, nil
}

func (ir ItemRepo) FetchItemById(id int) (models.Item, error) {
	openlog.Debug("fetching item by id.")
	item := models.Item{}
	err := ir.DB.First(&item, id).Error
	if err != nil {
		openlog.Error("error occured while fetching item by id")
		if err == gorm.ErrRecordNotFound {
			return item, common.ErrResourceNotFound
		}
		return item, errors.New("error occured while fetching item by id")
	}
	return item, nil
}

// Implement the FetchTotalCount method to satisfy the ItemRepositoryInterface
func (rr ItemRepo) FetchTotalCount() (int64, error) {

	openlog.Info("Fetching count of items from database")
	var count int64
	item := []models.Item{}
	result := rr.DB.Find(&item).Count(&count)
	if result.Error != nil {
		openlog.Error("Error occured while fetching count of items from database")
		return count, result.Error
	}
	return count, nil
}

func (ir ItemRepo) FetchAllItemsCount(filters string) (int64, error) {
	openlog.Debug("counting items.")
	var count int64
	filterScope := common.GetAllCondition(filters)
	err := ir.DB.Scopes(filterScope...).Model(&models.Item{}).Count(&count).Error
	if err != nil {
		openlog.Error("error occurred while counting items." + err.Error())
		return 0, errors.New("error occurred while counting items")
	}

	return count, nil
}
func (ir ItemRepo) Insert(item models.Item) (models.Item, error) {
	openlog.Debug("Inserting item into database")
	res := ir.DB.Create(&item)
	if res.Error != nil {
		openlog.Error("Error occured while inserting items into database")
		return item, res.Error
	}
	return item, nil
}

func (ir ItemRepo) GetItemByNameOrSKU(name, sku string) (models.Item, error) {
	var item models.Item
	err := ir.DB.Where("name = ? OR sku = ?", name, sku).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Item{}, nil // Item not found, return empty item with no error
		}
		return models.Item{}, err // Other error occurred
	}
	return item, nil // Item found, return the item with no error
}

func (ir ItemRepo) BulkInsertItems(items []models.Item) (int64, error) {
	openlog.Debug("Bulk inserting items")
	chunkSize := 100
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		res := ir.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sku"}}, // to ref unique sku column
			DoNothing: true,
			// UpdateAll: true,
		}).Create(items[i:end])
		if res.Error != nil {
			openlog.Error("Error occured while inserting the items")
			return 0, res.Error
		}
	}
	return int64(len(items)), nil
}

// get item by id from database
func (ir ItemRepo) GetById(id int) (models.Item, error) {
	openlog.Info("Fetching item by id from database")
	item := models.Item{}
	result := ir.DB.First(&item, id)
	//result := sr.DB.Model(&models.Skillset{}).Preload("Users").Where("id = ?", id).First(&skillset)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			openlog.Error("items not found")
			return item, common.ErrResourceNotFound
		}
		openlog.Error("Error occured while fetching item by id from database")
		return item, result.Error
	}

	return item, nil
}

// update item by id from database
func (ir ItemRepo) UpdateById(id int, item models.Item) (models.Item, error) {
	openlog.Info("Updating item by id from database")
	result := ir.DB.Model(&item).Where("id = ?", id).Updates(item)
	if result.Error != nil {
		openlog.Error("Error occured while updating items by id from database")
		return item, result.Error
	}
	return item, nil
}

// soft delete item by id from database
func (ir ItemRepo) DeleteById(id int) error {
	openlog.Info("Deleting item by id from database")
	result := ir.DB.Where("id = ?", id).Delete(&models.Item{})
	if result.Error != nil {
		openlog.Error("Error occured while deleting item by id from database")
		return result.Error
	}
	return nil
}
