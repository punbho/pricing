package plans

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Plan variables
var (
	PlanPath           = "plans.json"
	SupportedPlansPath = "supported_plans.json"
)

type PlanFileDao string

func (PlanFileDao) getSupportedPlans() ([]PlanDescription, error) {
	plans := []PlanDescription{}
	raw, err := ioutil.ReadFile(SupportedPlansPath)
	if err != nil {
		return []PlanDescription{}, errors.New("Unable to read supported plans file path")
	}

	err = json.Unmarshal(raw, &plans)
	if err != nil {
		return []PlanDescription{}, errors.New("Failed to unmarshal supported plans data")
	}
	return plans, nil
}

func (f PlanFileDao) getPlanDescription(planName string) (PlanDescription, error) {
	plans, err := f.getSupportedPlans()
	if err != nil {
		return PlanDescription{}, err
	}

	for _, plan := range plans {
		if plan.PlanName == planName {
			return plan, nil
		}
	}
	return PlanDescription{}, errors.New("Could not find information for this plan")
}

func (PlanFileDao) getAllPlanPrices() (CountryPlanMap, error) {
	var planMap CountryPlanMap
	raw, err := ioutil.ReadFile(PlanPath)
	if err != nil {
		return planMap, errors.New("Unable to read plans file path")
	}
	err = json.Unmarshal(raw, &planMap)
	if err != nil {
		return planMap, errors.New("Failed to unmarshal plans data")
	}
	return planMap, nil
}

func (f PlanFileDao) getPlanPricesForCountry(countryCode string) (CountryPlan, error) {
	planMap, err := f.getAllPlanPrices()
	if err != nil {
		return CountryPlan{}, err
	}
	if _, ok := planMap[countryCode]; !ok {
		return CountryPlan{}, errors.New("Unable to find the country with the given code")
	}
	return planMap[countryCode], nil
}

func (f PlanFileDao) getIndividualPlanPriceForCountry(countryCode string, planName string) (Plan, error) {
	countryPlanPrices, err := f.getPlanPricesForCountry(countryCode)
	if err != nil {
		return Plan{}, err
	}
	if _, ok := countryPlanPrices.PlanPrice[planName]; !ok {
		return Plan{}, errors.New("Plan does not exist in this country")
	}
	return countryPlanPrices.PlanPrice[planName], nil
}

func (PlanFileDao) updateAllPlanPrices(countryPlanPriceMap CountryPlanMap) error {
	planBytes, err := json.Marshal(countryPlanPriceMap)
	err = ioutil.WriteFile(PlanPath, planBytes, 0755)
	if err != nil {
		return errors.New("Failed to update the plan prices")
	}
	return nil
}

func (f PlanFileDao) updatePlanPricesForCountry(countryCode string, planPriceDetails map[string]Plan) error {
	planMap, err := f.getAllPlanPrices()
	if err != nil {
		return err
	}
	countryMap := CountryPlan{
		Currency:  planMap[countryCode].Currency,
		PlanPrice: planPriceDetails,
	}
	planMap[countryCode] = countryMap
	return f.updateAllPlanPrices(planMap)
}

func (PlanFileDao) isValidPlan(planName string, plan_info []PlanDescription) bool {
	for _, plan := range plan_info {
		if plan.PlanName == planName {
			return true
		}
	}
	return false
}
