package models

// See https://docs.goshippo.com/api-reference/parcels/create-a-new-parcel
type ParcelInput struct {
	Length       string       `json:"length,omitempty"`
	Width        string       `json:"width,omitempty"`
	Height       string       `json:"height,omitempty"`
	DistanceUnit string       `json:"distance_unit,omitempty"`
	Weight       string       `json:"weight"`
	MassUnit     string       `json:"mass_unit"`
	Template     string       `json:"template,omitempty"` // https://docs.goshippo.com/api-reference/parcel-templates/index
	Extra        *ParcelExtra `json:"extra,omitempty"`
	Metadata     string       `json:"metadata,omitempty"`
}

type ParcelExtra struct {
	COD        *ParcelCOD       `json:"COD,omitempty"`
	Insurance  *ParcelInsurance `json:"insurance,omitempty"`
	Reference1 string           `json:"reference_1,omitempty"`
	Reference2 string           `json:"reference_2,omitempty"`
}

type ParcelCOD struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

type ParcelInsurance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Provider string `json:"provider,omitempty"` // Omit for Shippo's third-party insurance.
	Content  string `json:"content"`
}

// See https://docs.goshippo.com/api-reference/parcels/retrieve-an-existing-parcel
type Parcel struct {
	ParcelInput
	CommonOutputFields
	ObjectState string `json:"object_state,omitempty"`
	Test        bool   `json:"test"`
}
