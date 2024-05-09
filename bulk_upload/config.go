package bulk_upload

import (
	"lp_customer_portal/database"
	"lp_customer_portal/models"
	"lp_customer_portal/repository/customer_repo"
	"lp_customer_portal/repository/item_repo"
	"lp_customer_portal/repository/order_repo"

	//"lp_customer_portal/repository/vehicle_repo"
	//zone_repository "lp_oms/repository/zone_repo"

	"gorm.io/gorm"
)

type BulkInsertConfig struct {
	AutoInsertCustomer bool
	AutoInsertAddress  bool
	AutoInsertItems    bool
	AutoInsertVehicles bool
	DeterminzeLatLng   bool // if true, then use google map api to determine the lat and lng
	DertermineZone     bool
	RejectDuplicate    bool // we will not insert duplicate orders
	TemplateName       string
}

type BulkInsert struct {
	Config       BulkInsertConfig
	DB           *gorm.DB
	OrderRepo    order_repo.OrderRepoInterface
	CustomerRepo customer_repo.CustomerRepositoryInterface
	//ZoneRepo     zone_repository.ZoneRepositoryInterface
	ItemRepo item_repo.ItemRepositoryInterface
	//	VehicleRepo  vehicle_repo.VehicleRepoInterface
	OrderRows   []OrderRow
	ItemRows    []ItemRow
	VehicleRows []VehicleRow
	Items       []models.Item     // these will be preloaded into system to avoid multiple db calls
	Customers   []models.Customer // these will be preloaded into system to avoid multiple db calls
	Zones       []models.Zone
	Addresses   []models.Address
	Vehicles    []models.Vehicle
}

func NewBulkInsert(config BulkInsertConfig) BulkInsert {
	db := database.GetClient()
	// db = db.Begin()
	or := order_repo.OrderRepo{DB: db}
	cr := customer_repo.CustomerRepository{DB: db}
	i := &item_repo.ItemRepo{DB: db}
	//v := vehicle_repo.VehicleRepo{DB: db}
	//	z := zone_repository.ZoneRepo{DbClient: db}

	bulkInsert := BulkInsert{
		Config:       config,
		OrderRepo:    or,
		CustomerRepo: cr,
		ItemRepo:     i,
		//VehicleRepo:  v,
		//ZoneRepo:     z,
		DB: db,
	}
	return bulkInsert
}
