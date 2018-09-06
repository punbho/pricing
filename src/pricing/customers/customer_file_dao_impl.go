package customers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Customer file path
var (
	CustomerFilePath = "customers.json"
)

// CustomerFileDao is the type for file based implementation for customer data.
type CustomerFileDao string

func (CustomerFileDao) getAllCustomers() ([]Customer, error) {
	customers := []Customer{}
	raw, err := ioutil.ReadFile(CustomerFilePath)
	if err != nil {
		return []Customer{}, errors.New("Unable to read customer data")
	}

	err = json.Unmarshal(raw, &customers)
	if err != nil {
		return []Customer{}, errors.New("Failed to unmarshal customer data")
	}
	return customers, nil
}

func (c CustomerFileDao) getCustomersByCountry(countryCode string) ([]Customer, error) {
	if _, ok := customerMap[countryCode]; !ok {
		return []Customer{}, errors.New("No customer is present for the passed country")
	}
	return customerMap[countryCode], nil
}
