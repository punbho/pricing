package customers

// This map holds customer information based on country to easily get customers from a country
// This is equivalent to creating an index in a database
var (
	customerMap map[string][]Customer
)

// Customer struct holds the data for the customers. The price of the plan is not saved with the customer
// the price can be fetched from the plan data
type Customer struct {
	Email       string `json:"email"`
	CountryCode string `json:"countryCode"`
	PlanName    string `json:"planName"`
}

// CustomerResponse struct is used to return the data of the customers as API response including price.
type CustomerResponse struct {
	Customer
	Price uint `json:"price"`
}

// CustomerDao interface provides methods to get all customer data and data based on the country.
type CustomerDao interface {
	getAllCustomers() ([]Customer, error)
	getCustomersByCountry(countryCode string) ([]Customer, error)
}
