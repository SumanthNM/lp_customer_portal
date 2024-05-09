/**
 * Sample Chassis Handler to print log
 *
**/

package chassisHandlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"lp_customer_portal/common"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/openlog"
	"github.com/xeipuuv/gojsonschema"
)

const Name = "Paylaod-Validator"

type PayloadValidatorHanldlerHandler struct{}

func init() { handler.RegisterHandler(Name, New) }

func New() handler.Handler { return &PayloadValidatorHanldlerHandler{} }

func (h *PayloadValidatorHanldlerHandler) Name() string { return Name }

func (h *PayloadValidatorHanldlerHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	// request object
	var req *http.Request
	if r, ok := inv.Args.(*http.Request); ok {
		req = r
	} else if r, ok := inv.Args.(*restful.Request); ok {
		req = r.Request
	} else {
		openlog.Error(fmt.Sprintf("this handler only works for http protocol, wrong type: %t", inv.Args))
		return
	}
	payload_bytes, err := ioutil.ReadAll(req.Body)
	openlog.Debug("got request to " + inv.URLPath)
	schemaPath := getSchema(inv.URLPath)
	openlog.Info(schemaPath)
	if schemaPath == "" {
		chain.Next(inv, func(r *invocation.Response) {
			cb(r)
		})
		return
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	documentLoader := gojsonschema.NewBytesLoader(payload_bytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		openlog.Error("error occured here" + err.Error())
		data := common.HTTPResponse{Msg: "Invalid Json", Status: 400}
		cb(&invocation.Response{Err: errors.New("invalid Json"), Status: 400, Result: data})
		fmt.Println(data)
		return
	}
	if result.Valid() {
		openlog.Info("Payload Validation completed")
		chain.Next(inv, cb)
		return
	} else {

		validationErrors := make([]string, 0)
		for _, desc := range result.Errors() {
			// // fmt.Printf("-----> %s\n", desc)
			validationErrors = append(validationErrors, desc.String())
			// validationErrors += desc.String() + "\n"
		}

		data := common.HTTPResponse{Msg: "Invalid Json", Status: 400, Data: validationErrors}
		fmt.Println(data)
		cb(&invocation.Response{Err: errors.New("Invalid Json"), Status: 400, Result: data})
		openlog.Error("schema validation errors")
		// handler.WriteBackErr(errors.New(validationErrors), 400, cb)

		return
	}
}

func getSchema(uri string) string {
	switch uri {
	case "/users":
		openlog.Debug("Schema Identified as createUserPayload ")
		return getSchemaPathFromConfig("createUserPayload")
	default:
		return ""
	}
}

func getSchemaPathFromConfig(schema string) string {
	return archaius.GetString(schema, "")
}
