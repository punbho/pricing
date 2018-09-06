package plans

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	dao PlanDao
)

func InitializeHandler(planDao PlanDao) {
	dao = planDao
}

func GetAllPlansHandler(w http.ResponseWriter, r *http.Request) {
	var planMap CountryPlanMap
	planMap, err := dao.getAllPlanPrices()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(planMap)
}

func GetAllPlansService() (CountryPlanMap, error) {
	var planMap CountryPlanMap
	planMap, err := dao.getAllPlanPrices()
	if err != nil {
		return planMap, err
	}
	return planMap, nil
}

func GetPlansForCountryHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	countryCode := pathParams["countryCode"]
	countryPlan, err := dao.getPlanPricesForCountry(countryCode)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(countryPlan)
}

func GetIndividualPlanDetailsPerCountryHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	countryCode := pathParams["countryCode"]
	planName := pathParams["planName"]
	plan, err := dao.getIndividualPlanPriceForCountry(countryCode, planName)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	planInfo, err := dao.getPlanDescription(planName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := PlanResponse{
		planInfo,
		plan,
	}
	json.NewEncoder(w).Encode(resp)
}

func ChangePlanPricesPerCountryHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	countryCode := pathParams["countryCode"]
	countryPlan, err := dao.getPlanPricesForCountry(countryCode)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	planPriceDetails := countryPlan.PlanPrice
	plansDecoder := json.NewDecoder(r.Body)
	var newPlanDetails map[string]uint
	err = plansDecoder.Decode(&newPlanDetails)
	if err != nil {
		http.Error(w, "Invalid request body format", 400)
		return
	}
	supportedPlans, err := dao.getSupportedPlans()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	timestamp := time.Now().Unix()
	for k, v := range newPlanDetails {
		if !dao.isValidPlan(k, supportedPlans) {
			http.Error(w, "Invalid plan "+k, 400)
			return
		}
		plan := Plan{
			Price:     v,
			ValidFrom: timestamp,
		}
		if _, ok := planPriceDetails[k]; ok {
			plan.OldPrice = planPriceDetails[k].Price
		}
		planPriceDetails[k] = plan
	}

	err = dao.updatePlanPricesForCountry(countryCode, planPriceDetails)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}
