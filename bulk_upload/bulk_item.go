package bulk_upload

import (
	"errors"
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	"strconv"

	"github.com/go-chassis/openlog"
)

// LoadCSV reads the csv, users the key map to do data transformation into itemRow
func (b *BulkInsert) ProcessItemRecords(records [][]string) (int64, error) {
	openlog.Debug("Processing records")
	openlog.Debug("Starting the transaction")
	database.StartTransaction()

	// get headers
	headers := records[0]
	// remove header from remaining data
	records = records[1:]

	itemRows := make([]ItemRow, 0)
	uniqueSKU := make(map[string]bool)
	for _, record := range records {
		data := make(map[string]string)

		for i := 0; i < len(record); i++ {
			data[headers[i]] = record[i]
		}
		// deduplication
		keymap := GetKeyMap("item", b.Config.TemplateName, ITEMKEYS)
		itemSKU := data[keymap["itemSKU"]]

		if itemSKU == "" {
			continue
		}
		if _, exists := uniqueSKU[itemSKU]; !exists {
			uniqueSKU[itemSKU] = true
			item, err := transformItem(data, keymap)
			if err != nil {
				openlog.Error("Error while transforming the data")
				openlog.Debug("Rolling back the transaction")
				database.RollbackTransaction()
				return 0, err
			}
			itemRows = append(itemRows, *item)
		}
	}

	// convert ItemRows to Items and remove duplicates
	itemModelsList := []models.Item{}
	for _, itemRow := range itemRows {
		item := getItem(itemRow)
		itemModelsList = append(itemModelsList, item)
	}
	// upsert into DB
	count, err := b.ItemRepo.BulkInsertItems(itemModelsList)
	if err != nil {
		openlog.Error("Error occured while performing db bulk insert. [" + err.Error() + "]")
		openlog.Debug("rollback the transaction")
		database.RollbackTransaction() // rollback the transaction
		return 0, err
	}

	openlog.Debug("commiting the trasaction")
	database.CommitTransaction()
	openlog.Debug("Number of items inserted [" + strconv.FormatInt(count, 10) + "]")
	return count, nil
}

func transformItem(data map[string]string, keymap map[string]string) (*ItemRow, error) {
	item := ItemRow{}
	item.ItemSKU = data[keymap["itemSKU"]]
	item.Name = data[keymap["name"]]
	if len(data[keymap["weight"]]) == 0 {
		item.Weight = 10
	} else {
		w, err := strconv.ParseFloat(data[keymap["weight"]], 64)
		if err != nil {
			openlog.Error("Error while converting the weight to float [" + err.Error() + "]" + data[keymap["weight"]])
			return nil, errors.New("error occured while parsing weight")
		}
		item.Weight = w
	}

	if value, ok := data[keymap["innerPack"]]; ok && value != "" {
		innerPack, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			openlog.Error("Error while converting the innerpack to int [" + err.Error() + "]" + data[keymap["innerPack"]])
			return nil, errors.New("error occurred while parsing innerpack")
		}
		item.InnerPack = int(innerPack)
	}

	return &item, nil
}
