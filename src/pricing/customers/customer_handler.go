package customers

import (
	"encoding/json"
	"net/http"
	"pricing/plans"

	"github.com/gorilla/mux"
)

var (
	dao CustomerDao
)

// InitializeHandler initializes the dao for the customer data.
func InitializeHandler(custDao CustomerDao) {
	dao = custDao
}

func GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {
	customers, err := dao.getAllCustomers()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cutomerResponses, err := getCustomerResponses(customers)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(cutomerResponses)
}

func GetCustomersByCountryHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	countryCode := pathParams["countryCode"]
	customers, err := dao.getCustomersByCountry(countryCode)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cutomerResponses, err := getCustomerResponses(customers)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(cutomerResponses)
}

func PreprocessCustomerData() error {
	customers, err := dao.getAllCustomers()
	if err != nil {
		return err
	}
	customerMap = make(map[string][]Customer)
	for _, cust := range customers {
		if _, ok := customerMap[cust.CountryCode]; !ok {
			customerMap[cust.CountryCode] = []Customer{cust}
		} else {
			customerMap[cust.CountryCode] = append(customerMap[cust.CountryCode], cust)
		}
	}
	return nil
}

func getCustomerResponses(customers []Customer) ([]CustomerResponse, error) {
	plans, err := plans.GetAllPlansService()
	if err != nil {
		return []CustomerResponse{}, err
	}
	customerResponses := []CustomerResponse{}
	for _, cust := range customers {
		customerResponse := CustomerResponse{
			Customer: cust,
			Price:    plans[cust.CountryCode].PlanPrice[cust.PlanName].Price,
		}
		customerResponses = append(customerResponses, customerResponse)
	}
	return customerResponses, nil
}
