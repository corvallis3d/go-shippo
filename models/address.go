package models

const (
	// Deprecated: The current API no longer returns object_source.
	AddressObjectSourceFullyEntered = "FULLY_ENTERED"
	// Deprecated: The current API no longer returns object_source.
	AddressObjectSourcePartiallyEntered = "PARTIALLY_ENTERED"
	// Deprecated: The current API no longer returns object_source.
	AddressObjectSourceValidator = "VALIDATOR"
)

// See https://docs.goshippo.com/api-reference/addresses/create-a-new-address
type AddressInput struct {
	Name          string `json:"name,omitempty"`
	Company       string `json:"company,omitempty"`
	Street1       string `json:"street1,omitempty"`
	StreetNo      string `json:"street_no,omitempty"`
	Street2       string `json:"street2,omitempty"`
	Street3       string `json:"street3,omitempty"`
	City          string `json:"city,omitempty"`
	Zip           string `json:"zip,omitempty"`
	State         string `json:"state,omitempty"`
	Country       string `json:"country,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	IsResidential bool   `json:"is_residential"`
	Validate      bool   `json:"validate"`
	Metadata      string `json:"metadata,omitempty"`
}

// See https://docs.goshippo.com/api-reference/addresses/create-a-new-address
type Address struct {
	AddressInput
	CommonOutputFields
	IsComplete        bool               `json:"is_complete"`
	Test              bool               `json:"test"`
	ValidationResults *ValidationResults `json:"validation_results"`
}

// See https://docs.goshippo.com/api-reference/addresses/validate-an-address
type ValidationResults struct {
	IsValid  bool             `json:"is_valid"`
	Messages []*OutputMessage `json:"messages"`
}
