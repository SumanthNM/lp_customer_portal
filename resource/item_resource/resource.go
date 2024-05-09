package item_resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"lp_customer_portal/bulk_upload"
	"lp_customer_portal/common"
	"lp_customer_portal/schemas"
	"lp_customer_portal/services/item_services"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
)

type ItemResource struct {
}

func (ir *ItemResource) Inject(service item_services.ItemServiceInterface) {

}

func (ir *ItemResource) FetchAllItems(context *restful.Context) {
	openlog.Debug("Got a request to fetch all items")
	pageno, limit := common.GetPaginationParams(context.ReadQueryParameter("pageno"), context.ReadQueryParameter("limit"))
	filters := context.ReadQueryParameter("filters")
	// get all roles
	is := item_services.New()
	res := is.FetchAllItems(pageno, limit, filters)
	context.WriteJSON(res, "application/json")
}

func (ir *ItemResource) FetchItemById(context *restful.Context) {
	// read path params
	openlog.Debug("Got a request to fetch item by id")
	idStr := context.ReadPathParameter("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		openlog.Error(err.Error())
		context.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Could not read payload"}, "application/json")
		return
	}
	is := item_services.New()
	res := is.FetchItemById(id)
	context.WriteJSON(res, "application/json")
}

func (ir *ItemResource) BulkInsertItems(ctx *restful.Context) {
	openlog.Info("Received request to insert ")
	startTime := time.Now()
	file_data := make([]byte, 0)
	file_name := ""
	config := bulk_upload.BulkInsertConfig{}
	// read mulitpart form data.
	multipart, err := ctx.ReadRequest().MultipartReader()
	if err != nil {
		openlog.Error("Error occured while reading the multipart data.")
		ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
		return
	}
	for {
		part, err := multipart.NextPart()
		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "file":
			file_name = part.FileName()
			file_data, err = ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		case "templateType":
			ac, err := ioutil.ReadAll(part)
			if err != nil {
				openlog.Error("Error occured while reading file")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
			config.TemplateName = string(ac[:])
			if err != nil {
				openlog.Error("Error occured while reading payload")
				ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
				return
			}
		}
		if file_name == "" {
			openlog.Error("Error occured while reading the multipart data.")
			ctx.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Error occured while reading payload"}, "application/json")
			return
		}
	}
	is := item_services.New()
	res := is.BulkInsertItems(file_data, file_name, config)
	endTime := time.Since(startTime)
	fmt.Printf("Execution time: %v\n", endTime)
	ctx.WriteHeaderAndJSON(res.Status, res, "application/json")

}

// creates the item
func (ir *ItemResource) CreateItem(context *restful.Context) {
	openlog.Info("Got a request to create item")
	item := schemas.ItemPayload{}

	// Read the payload into item from context
	err := context.ReadEntity(&item)
	if err != nil {
		openlog.Error(err.Error())
		// Send error response
		context.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Could not read paylaod"}, "application/json")
		return
	}
	// validation
	payloadErrors := common.ValidateStruct(item)
	if payloadErrors != nil {
		openlog.Error("Validation Errors")
		context.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Validation Errors", Data: payloadErrors}, "application/json")
		return
	}
	// call service layer
	is := item_services.New()
	res := is.CreateItem(item)
	context.WriteHeaderAndJSON(res.Status, res, "application/json")
}

func (ir *ItemResource) UpdateItemById(context *restful.Context) {
	// read path params
	id, _ := strconv.Atoi(context.ReadPathParameter("id"))

	item := schemas.ItemPayload{}
	err := context.ReadEntity(&item)
	if err != nil {
		openlog.Error(err.Error())
		context.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Could not read payload"}, "application/json")
		return
	}

	// validation
	payloadErrors := common.ValidateStruct(item)
	if payloadErrors != nil {
		openlog.Error("Validation Errors")
		context.WriteHeaderAndJSON(400, common.HTTPResponse{Status: 400, Msg: "Validation Errors", Data: payloadErrors}, "application/json")
		return
	}

	// update item by id
	is := item_services.New()
	res := is.UpdateItemById(id, item)
	context.WriteJSON(res, "application/json")
}

func (ir *ItemResource) DeleteItemById(context *restful.Context) {
	// read path params
	id, _ := strconv.Atoi(context.ReadPathParameter("id"))
	// delete vehicle by id
	is := item_services.New()
	res := is.DeleteItemById(id)
	context.WriteJSON(res, "application/json")
}

// Define all APIs here.
func (r *ItemResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: common.Basepath + "/customer/items", ResourceFunc: r.FetchAllItems, Consumes: []string{"application/json"}, Produces: []string{"application/json"}},
		{Method: http.MethodGet, Path: common.Basepath + "/customer/items/{id}", ResourceFunc: r.FetchItemById, Consumes: []string{"application/json"}, Produces: []string{"application/json"}},
		{Method: http.MethodDelete, Path: common.Basepath + "/customer/items/{id}", ResourceFunc: r.DeleteItemById, Consumes: []string{"application/json"}, Produces: []string{"application/json"}},
		{Method: http.MethodPost, Path: common.Basepath + "/customer/items", ResourceFunc: r.CreateItem, Consumes: []string{"application/json"}, Produces: []string{"application/json"}},
		{Method: http.MethodPut, Path: common.Basepath + "/customer/items/{id}", ResourceFunc: r.UpdateItemById, Consumes: []string{"application/json"}, Produces: []string{"application/json"}},

		{Method: http.MethodPost, Path: common.Basepath + "/items/bulk-insert", ResourceFunc: r.BulkInsertItems, Consumes: []string{"multipart/form-data"}, Produces: common.API_HEADERS},
	}
}
