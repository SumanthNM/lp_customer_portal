package here_helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"lp_customer_portal/common"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
)

const url = "https://geocode.search.hereapi.com/v1/geocode"
const api_key_q = "apiKey"
const q = "q"

func GetCoordinates(address string) (float64, float64, error) {
	// recover
	defer func() {
		if r := recover(); r != nil {
			openlog.Error("Error getting coordinates for address: " + address)
		}
	}()

	openlog.Debug("Getting coordinates for address: " + address)
	// create request
	api_key := archaius.GetString("here.apiKey", "")
	openlog.Debug("Here API key: " + api_key)
	headers := map[string]string{}
	queryParms := map[string]string{}
	queryParms[api_key_q] = api_key
	queryParms[q] = address
	res, err := common.GETReqeust(url, queryParms, headers)
	if err != nil {
		openlog.Error("Error getting coordinates for address: " + fmt.Sprint(url))
		return -1, -1, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		openlog.Error("Error getting coordinates for address: " + address)
		openlog.Error(err.Error())
		return -1, -1, err
	}
	// check  if response has items
	if _, ok := data["items"]; !ok {
		openlog.Error("Error getting coordinates for address: " + address)
		return -1, -1, err
	}
	items := data["items"].([]interface{})
	if len(items) == 0 {
		openlog.Error("Error getting coordinates for address: " + address)
		return -1, -1, errors.New("coordinates not found for the address")
	}

	position := items[0].(map[string]interface{})["position"].(map[string]interface{})
	lat := position["lat"].(float64)
	lng := position["lng"].(float64)
	return lat, lng, nil
}
