package models

import "encoding/json"

// See https://docs.goshippo.com/api-reference/batches/create-a-batch
type BatchInput struct {
	DefaultCarrierAccount    string                `json:"default_carrier_account,omitempty"`
	DefaultServiceLevelToken string                `json:"default_servicelevel_token,omitempty"`
	LabelFileType            string                `json:"label_filetype,omitempty"`
	Metadata                 string                `json:"metadata,omitempty"`
	BatchShipments           []*BatchShipmentInput `json:"batch_shipments"`
}

// See https://docs.goshippo.com/api-reference/batches/retrieve-a-batch
type Batch struct {
	CommonOutputFields
	Status                   string              `json:"status,omitempty"`
	DefaultCarrierAccount    string              `json:"default_carrier_account,omitempty"`
	DefaultServiceLevelToken string              `json:"default_servicelevel_token,omitempty"`
	LabelFileType            string              `json:"label_filetype,omitempty"`
	Metadata                 string              `json:"metadata,omitempty"`
	BatchShipments           *BatchShipmentList  `json:"batch_shipments"`
	LabelURL                 []string            `json:"label_url"`
	ObjectResults            *BatchObjectResults `json:"object_results"`
	Test                     bool                `json:"test"`
}

type BatchShipmentList struct {
	Count    int              `json:"count"`
	Next     string           `json:"next,omitempty"`
	Previous string           `json:"previous,omitempty"`
	Results  []*BatchShipment `json:"results"`
}

type BatchObjectResults struct {
	PurchaseSucceeded int `json:"purchase_succeeded"`
	PurchaseFailed    int `json:"purchase_failed"`
	CreationFailed    int `json:"creation_failed"`
	CreationSucceeded int `json:"creation_succeeded"`
}

// See https://docs.goshippo.com/api-reference/batches/add-shipments-to-a-batch
type BatchShipmentInput struct {
	CarrierAccount    string `json:"carrier_account,omitempty"`
	ServiceLevelToken string `json:"servicelevel_token,omitempty"`
	Shipment          string `json:"shipment"`
	Metadata          string `json:"metadata,omitempty"`
	// ShipmentDetails supplies the nested shipment required by the current API.
	// When set, it takes precedence over the legacy Shipment ID.
	ShipmentDetails *ShipmentInput `json:"-"`
}

// See https://docs.goshippo.com/api-reference/batches/retrieve-a-batch
type BatchShipment struct {
	BatchShipmentInput
	CommonOutputFields
	Status      string            `json:"status,omitempty"`
	Transaction string            `json:"transaction"`
	Messages    []json.RawMessage `json:"messages"`
}

// MarshalJSON supports a nested shipment while preserving the legacy ID field.
func (s BatchShipmentInput) MarshalJSON() ([]byte, error) {
	type batchShipmentInput BatchShipmentInput
	var shipment interface{} = s.Shipment
	if s.ShipmentDetails != nil {
		shipment = s.ShipmentDetails
	}
	return json.Marshal(struct {
		batchShipmentInput
		Shipment interface{} `json:"shipment"`
	}{batchShipmentInput: batchShipmentInput(s), Shipment: shipment})
}

// MarshalJSON keeps response fields when the embedded request type implements
// json.Marshaler. Batch shipment responses contain a shipment ID, not a request.
func (s BatchShipment) MarshalJSON() ([]byte, error) {
	type batchShipmentInput BatchShipmentInput
	return json.Marshal(struct {
		batchShipmentInput
		CommonOutputFields
		Status      string            `json:"status,omitempty"`
		Transaction string            `json:"transaction"`
		Messages    []json.RawMessage `json:"messages"`
	}{batchShipmentInput(s.BatchShipmentInput), s.CommonOutputFields, s.Status, s.Transaction, s.Messages})
}
