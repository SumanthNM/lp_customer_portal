/**
 * Contains the constants requried for HTTP methods
 * Add all Response messages here as constants.
**/

package common

type HTTPResponse struct {
	Msg    string      `json:"_msg"`
	Status int         `json:"_status"`
	Data   interface{} `json:"data"`
}

var API_HEADERS = []string{"application/json"}

var JSON_HEADER = "application/json"
var ORDER_TIME_FORMAT = "02/01/2006" //"02/01/2006"

var Basepath = "api/v1"
var TIME_FORMAT = "2006-01-02 15:04:05"

var DATE_FORMAT = "2006-01-02"
