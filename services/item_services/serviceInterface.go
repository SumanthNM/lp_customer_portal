/**
 * Inteface for service layer
 *
**/

package item_services

import (
	"lp_customer_portal/bulk_upload"
	common "lp_customer_portal/common"
	"lp_customer_portal/schemas"
)

type ItemServiceInterface interface {
	FetchAllItems(pageno, limit int, filters string) common.HTTPResponse
	FetchItemById(id int) common.HTTPResponse
	BulkInsertItems(fileData []byte, filename string, config bulk_upload.BulkInsertConfig) common.HTTPResponse
	CreateItem(item schemas.ItemPayload) common.HTTPResponse
	UpdateItemById(id int, item schemas.ItemPayload) common.HTTPResponse
	DeleteItemById(id int) common.HTTPResponse
}
