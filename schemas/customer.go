package schemas

import "lp_customer_portal/models"

type Address struct {
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	ZipCode   uint    `json:"zipCode"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ZoneID    uint    `json:"zoneId"`
}

func (a *Address) ToModel(model *models.Address) error {
	model.Street = a.Street
	model.City = a.City
	model.State = a.State
	model.Country = a.Country
	model.ZipCode = a.ZipCode
	model.Latitude = a.Latitude
	model.Longitude = a.Longitude
	return nil
}

type CustomerPayload struct {
	//CustomerID   string    `json:"customerId"`
	CustomerName string    `json:"customerName"`
	Contact      string    `json:"contact"`
	Email        string    `json:"email"`
	Address      []Address `json:"address"`
}

func (c *CustomerPayload) ToModel(model *models.Customer) error {
	//model.CustomerID = c.CustomerID
	model.CustomerName = c.CustomerName
	model.Contact = c.Contact
	model.Email = c.Email
	return nil
}
