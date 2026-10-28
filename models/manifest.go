package models

import "time"

// See https://docs.goshippo.com/api-reference/manifests/create-a-new-manifest
type ManifestInput struct {
	CarrierAccount string    `json:"carrier_account,omitempty"`
	ShipmentDate   time.Time `json:"shipment_date,omitempty"`
	AddressFrom    string    `json:"address_from,omitempty"`
	Transactions   []string  `json:"transactions,omitempty"`
	Documents      []string  `json:"documents,omitempty"`
	Async          bool      `json:"async"`
}

// See https://docs.goshippo.com/api-reference/manifests/retrieve-a-manifest
type Manifest struct {
	ManifestInput
	CommonOutputFields
	// Deprecated: Use Status for the current API.
	State  string   `json:"state,omitempty"`
	Status string   `json:"status,omitempty"`
	Errors []string `json:"errors,omitempty"`
}
