package models

import "time"

const (
	TrackingStatusStatusUnknown    = "UNKNOWN"
	TrackingStatusStatusPreTransit = "PRE_TRANSIT"
	TrackingStatusStatusDelivered  = "DELIVERED"
	TrackingStatusStatusTransit    = "TRANSIT"
	TrackingStatusStatusFailure    = "FAILURE"
	TrackingStatusStatusReturned   = "RETURNED"
)

// See https://goshippo.com/docs/reference#tracks
type TrackingStatusInput struct {
	Carrier               string `json:"carrier"`
	TrackingNumber        string `json:"tracking_number"`
	Metadata              string `json:"metadata,omitempty"`
	IncludePackageDetails bool   `json:"include_package_details,omitempty"`
}

// See https://goshippo.com/docs/reference#tracks
type TrackingStatus struct {
	TrackingStatusInput
	AddressFrom     *TrackingStatusLocation `json:"address_from,omitempty"`
	AddressTo       *TrackingStatusLocation `json:"address_to,omitempty"`
	ETA             time.Time               `json:"eta"`
	OriginalETA     time.Time               `json:"original_eta,omitempty"`
	Transaction     string                  `json:"transaction,omitempty"`
	Messages        []string                `json:"messages,omitempty"`
	Weight          []TrackingWeight        `json:"weight,omitempty"`
	Dimensions      []TrackingDimensions    `json:"dimensions,omitempty"`
	ServiceLevel    *ServiceLevel           `json:"servicelevel,omitempty"`
	TrackingStatus  *TrackingStatusDict     `json:"tracking_status,omitempty"`
	TrackingHistory []*TrackingStatusDict   `json:"tracking_history,omitempty"`
}

// TrackingWeight is a carrier-reported weight. Units and numeric values are
// returned as strings and may differ from the parcel used to purchase a label.
type TrackingWeight struct {
	Value string `json:"value,omitempty"`
	Unit  string `json:"unit,omitempty"`
}

// TrackingDimensions may contain partial dimensions supplied by the carrier.
type TrackingDimensions struct {
	Length string `json:"length,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

type TrackingStatusDict struct {
	CommonOutputFields
	Substatus     *TrackingStatusSubstatus `json:"substatus,omitempty"`
	Status        string                   `json:"status,omitempty"`
	StatusDetails string                   `json:"status_details,omitempty"`
	StatusDate    time.Time                `json:"status_date,omitempty"`
	Location      *TrackingStatusLocation  `json:"location,omitempty"`
}

type TrackingStatusSubstatus struct {
	Code           string `json:"code"`
	Text           string `json:"text"`
	ActionRequired bool   `json:"action_required"`
}

type TrackingStatusLocation struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}
