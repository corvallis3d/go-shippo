package models

// See https://docs.goshippo.com/api-reference/rates/retrieve-a-rate
type Rate struct {
	CommonOutputFields
	Attributes       []string      `json:"attributes"`
	AmountLocal      string        `json:"amount_local"`
	CurrencyLocal    string        `json:"currency_local"`
	Amount           string        `json:"amount"`
	Currency         string        `json:"currency"`
	Provider         string        `json:"provider"`
	ProviderImage75  string        `json:"provider_image_75"`
	ProviderImage200 string        `json:"provider_image_200"`
	ServiceLevel     *ServiceLevel `json:"servicelevel"`
	// Deprecated: Use EstimatedDays for the current API.
	Days                   int              `json:"days"`
	EstimatedDays          int              `json:"estimated_days"`
	Shipment               string           `json:"shipment"`
	ArrivesBy              string           `json:"arrives_by,omitempty"`
	IncludedInsurancePrice string           `json:"included_insurance_price,omitempty"`
	DurationTerms          string           `json:"duration_terms"`
	CarrierAccount         string           `json:"carrier_account"`
	Zone                   string           `json:"zone"`
	Messages               []*OutputMessage `json:"messages"`
	Test                   bool             `json:"test"`
}

type ServiceLevel struct {
	Name               string        `json:"name"`
	Token              string        `json:"token"` // https://docs.goshippo.com/api-reference/service-levels/index
	Terms              string        `json:"terms"`
	ExtendedToken      string        `json:"extended_token,omitempty"`
	ParentServiceLevel *ServiceLevel `json:"parent_servicelevel,omitempty"`
}
