/**
 * Server starts here
 * for reference
 * - openlog - https://github.com/go-chassis/openlog/blob/master/openlog.go
 * - Open architecture - https://medium.com/hackernoon/golang-clean-archithecture-efd6d7c43047
**/

package main

import (
	_ "lp_customer_portal/chassisHandlers"
	"lp_customer_portal/database"
	"lp_customer_portal/repository/customer_repo"
	"lp_customer_portal/repository/item_repo"
	"lp_customer_portal/repository/order_repo"

	//	"lp_customer_portal/repository/planner_repo"

	//"lp_customer_portal/repository/vehicle_repo"
	//zone_repository "lp_customer_portal/repository/zone_repo"
	"lp_customer_portal/resource/customer_resource"
	"lp_customer_portal/resource/item_resource"

	//optimizer "lp_customer_portal/resource/optimizer"
	"lp_customer_portal/resource/order_resource"

	//"lp_customer_portal/resource/vehicle_resource"
	"lp_customer_portal/services/customer_services"
	"lp_customer_portal/services/item_services"

	//optimizer_service "lp_customer_portal/services/optimizer"
	"lp_customer_portal/services/order_service"

	//vehicle_service "lp_customer_portal/services/vehicle_services"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
)

func getService() (*customer_services.CustomerService, *order_service.OrderService, *item_services.ItemService) {
	db := database.GetClient()
	customer_repo := &customer_repo.CustomerRepository{DB: db}
	order_repo := &order_repo.OrderRepo{DB: db}
	item_repo := &item_repo.ItemRepo{DB: db}
	customer_service := &customer_services.CustomerService{Repo: customer_repo}
	return customer_service,
		&order_service.OrderService{OrderRepo: order_repo, CustomerRepo: customer_repo, CustomerService: customer_service, ItemRepo: item_repo},
		&item_services.ItemService{ItemRepo: item_repo}

}

func main() {

	// Register schema
	customer_apis := customer_resource.CustomerResource{}
	order_apis := order_resource.OrderResource{}
	item_apis := item_resource.ItemResource{}

	chassis.RegisterSchema("rest", &customer_apis)
	chassis.RegisterSchema("rest", &order_apis)
	chassis.RegisterSchema("rest", &item_apis)

	if err := chassis.Init(); err != nil {
		openlog.Fatal("Init failed." + err.Error())
		return
	}

	// Add database configurations to archaius
	if err := archaius.AddFile("./conf/database.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}
	// Add orders paths configurations to archaius
	if err := archaius.AddFile("./conf/bulkupload.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}

	// Add orders paths here to archaius
	if err := archaius.AddFile("./conf/here.yaml"); err != nil {
		openlog.Error("add props here failed." + err.Error())
		return
	}

	// Add schema paths configurations to archaius
	if err := archaius.AddFile("./conf/payloadSchemas.yaml"); err != nil {
		openlog.Error("add props configurations failed." + err.Error())
		return
	}

	// Add schema paths configurations to archaius
	if err := archaius.AddFile("./conf/aws.yaml"); err != nil {
		openlog.Error("add aws configurations failed." + err.Error())
		return
	}

	// Add schema paths configurations to archaius
	if err := archaius.AddFile("./conf/vroom.yaml"); err != nil {
		openlog.Error("add aws configurations failed." + err.Error())
		return
	}

	// Add schema paths configurations to archaius
	if err := archaius.AddFile("./conf/versafleet.yaml"); err != nil {
		openlog.Error("add aws configurations failed." + err.Error())
		return
	}

	openlog.Debug("Connecting to the database.")
	// Server will not start if error occurs.
	if err := database.Connect(); err != nil {
		openlog.Fatal("Error occured while connecting to database")
		return
	}

	// Inject service into resource
	cs, os, is := getService()
	customer_apis.Inject(cs)
	order_apis.Inject(os)
	item_apis.Inject(is)

	chassis.Run()
}
