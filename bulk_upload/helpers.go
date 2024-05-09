package bulk_upload

import (
	"bytes"
	"encoding/csv"
	"lp_customer_portal/models"
	"strconv"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
	"github.com/xuri/excelize/v2"
)

func ReadCSV(data []byte) ([][]string, error) {
	// fix issue w BOM UTF-8 encoding that comes with thai character
	reader := csv.NewReader(utfbom.SkipOnly(bytes.NewReader(data)))
	records, err := reader.ReadAll()
	if err != nil {
		openlog.Error("Error while reading the csv file")
		return [][]string{}, err
	}
	openlog.Debug("Successfully read the csv ")
	openlog.Debug("Number of rows read [" + strconv.Itoa(len(records)) + "]")

	return records, nil
}

func ReadExcel(data []byte) ([][]string, error) {
	excelFile, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		openlog.Error("Error while reading the excel file")
		return [][]string{}, err
	}

	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		openlog.Error("Error while reading the excel file")
		return [][]string{}, err
	}
	openlog.Debug("Successfully read the excel ")
	openlog.Debug("Number of rows read [" + strconv.Itoa(len(rows)) + "]")

	return rows, nil
}

func GetFileType(file_name string) string {
	lowercaseFileName := strings.ToLower(file_name)

	if strings.HasSuffix(lowercaseFileName, ".csv") {
		return "csv"
	} else if strings.HasSuffix(lowercaseFileName, ".xlsx") {
		return "xlsx"
	} else {
		return "unknown"
	}
}

func GetKeyMap(uploadType, templateName string, keys []string) map[string]string {
	// read from archaius
	keymap := make(map[string]string)
	for _, v := range keys {
		keymap[v] = archaius.GetString(uploadType+"."+templateName+"."+v, v)
	}
	return keymap
}

func getItem(row ItemRow) models.Item {
	itemObj := models.Item{
		Name:      row.Name,
		SKU:       row.ItemSKU,
		Weight:    row.Weight,
		Volume:    row.Volume,
		InnerPack: row.InnerPack,
	}

	// to handle empty inner pack
	if row.InnerPack != 0 {
		itemObj.InnerPack = row.InnerPack
	} else {
		itemObj.InnerPack = 1
	}

	return itemObj
}

func parseAndFormatDate(dateString string, dateFormats []string) (time.Time, error) {
	var parsedTime time.Time
	var err error

	for _, format := range dateFormats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// required keys from bulkupload.yaml
var ITEMKEYS = []string{"itemSKU", "status", "name", "weight", "innerPack"}
var ORDERKEYS = []string{"orderNo", "customerName", "channel", "customerType", "itemName", "itemSKU", "quantityMain", "quantityMinor", "orderedDate", "zip", "deliveryDate", "deliveryStreetAddress"}
var A21ORDERKEYS = []string{"trackingNumber", "jobType", "customerName", "pickupStreetAddress", "pickupPostal", "pickupLatLng", "deliveryStreetAddress", "deliveryPostal", "deliveryLatLng", "orderedDate", "preferredPickupStartDateTime", "preferredPickupEndDateTime", "preferredDeliveryStartDateTime", "preferredDeliveryEndDateTime", "pickupCustomServiceTime", "deliveryCustomServiceTime", "zone", "skills", "quantity"}
var VEHICLEKEYS = []string{"licensePlate", "vehicleType", "depotStartId", "depotEndId"}

// list of date formats
var EXCEL_DATE_FORMATS = []string{
	"2006-01-02 15:04:05",
	"2/1/06 15:04",
}
