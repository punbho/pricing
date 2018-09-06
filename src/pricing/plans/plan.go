package plans

// Plan has the details about a specific plan including the prices and validity
type Plan struct {
	Price     uint  `json:"price"`
	ValidFrom int64 `json:"valid_from"`
	OldPrice  uint  `json:"old_price"`
}

// PlanResponse is the struct used to send individual plan information response
type PlanResponse struct {
	PlanDescription
	Plan
}

type PlanDescription struct {
	PlanName string `json:"plan_name"`
	PlanInfo string `json:"plan_info"`
}

// CountryPlan has the details specific to a particular country and the plans corresponding to that country.
type CountryPlan struct {
	Currency  string          `json:"currency"`
	PlanPrice map[string]Plan `json:"plan_price"` // Map of plan name to Plan struct
}

// CountryPlanMap is the overall map of country code to the CountryPlan struct
type CountryPlanMap map[string]CountryPlan

// PlanDao interface is the interface to get/set plan data from/to disk. Multiple implementations can be possible like using File, DB, etc.
type PlanDao interface {
	getSupportedPlans() ([]PlanDescription, error)
	getPlanDescription(planName string) (PlanDescription, error)
	getAllPlanPrices() (CountryPlanMap, error)
	getPlanPricesForCountry(countryCode string) (CountryPlan, error)
	getIndividualPlanPriceForCountry(countryCode string, planName string) (Plan, error)
	updatePlanPricesForCountry(countryCode string, planPriceDetails map[string]Plan) error
	updateAllPlanPrices(countryPlanPriceMap CountryPlanMap) error
	isValidPlan(planName string, plan_info []PlanDescription) bool
}
