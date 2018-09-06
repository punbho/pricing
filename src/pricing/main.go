package main

import (
	"fmt"
	"net/http"

	"pricing/customers"
	"pricing/plans"

	"github.com/gorilla/mux"
)

var (
	baseURI            = "/pricing/api/v1"
	allPlansURI        = baseURI + "/plans"
	plansPerCountryURI = baseURI + "/plans/{countryCode}"
	individualPlanURI  = baseURI + "/plans/{countryCode}/{planName}"

	allCustomersURI        = baseURI + "/customers"
	customersPerCountryURI = baseURI + "/customers/{countryCode}"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc(allPlansURI, plans.GetAllPlansHandler).Methods("GET")
	router.HandleFunc(plansPerCountryURI, plans.GetPlansForCountryHandler).Methods("GET")
	router.HandleFunc(plansPerCountryURI, plans.ChangePlanPricesPerCountryHandler).Methods("PUT")
	router.HandleFunc(individualPlanURI, plans.GetIndividualPlanDetailsPerCountryHandler).Methods("GET")

	router.HandleFunc(allCustomersURI, customers.GetAllCustomersHandler).Methods("GET")
	router.HandleFunc(customersPerCountryURI, customers.GetCustomersByCountryHandler).Methods("GET")
	http.Handle("/", router)

	var fileDao plans.PlanFileDao
	plans.InitializeHandler(fileDao)

	var customerDao customers.CustomerFileDao
	customers.InitializeHandler(customerDao)
	err := customers.PreprocessCustomerData()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server with error " + err.Error())
	}
}
