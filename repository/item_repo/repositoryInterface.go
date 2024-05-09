package item_repo

import "lp_customer_portal/models"

type ItemRepositoryInterface interface {
	Insert(item models.Item) (models.Item, error)
	GetItemByNameOrSKU(name, sku string) (models.Item, error)
	FetchAllItems(pageno, limit int, filters string) ([]models.Item, error)
	FetchItemById(id int) (models.Item, error)
	FetchTotalCount() (int64, error)
	BulkInsertItems(items []models.Item) (int64, error)
	GetById(id int) (models.Item, error)
	UpdateById(id int, item models.Item) (models.Item, error)
	DeleteById(id int) error
	FetchAllItemsCount(filters string) (int64, error)
}
