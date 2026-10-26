package models

import (
	"encoding/json"
	"time"
)

// See https://goshippo.com/docs/reference#transactions
type TransactionInput struct {
	Rate          string `json:"rate,omitempty"`
	Order         string `json:"order,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	LabelFileType string `json:"label_file_type,omitempty"`
	Async         bool   `json:"async"`

	Shipment       *ShipmentInput `json:"shipment,omitempty"`
	CarrierAccount string         `json:"carrier_account,omitempty"`
	// ServerLevelToken selects the service for an instant purchase. The historical
	// Go field name is retained for compatibility; the JSON key is servicelevel_token.
	ServerLevelToken string `json:"servicelevel_token,omitempty"`
}

// See https://goshippo.com/docs/reference#transactions
type Transaction struct {
	TransactionInput
	CommonOutputFields
	ObjectState string `json:"object_state,omitempty"`
	Status      string `json:"status,omitempty"`
	Test        bool   `json:"test"`
	Parcel      string `json:"parcel,omitempty"`
	// RateDetails is populated when Shippo returns an expanded rate (including
	// instant purchases). Rate still contains its object ID for existing callers.
	RateDetails          *CoreRate        `json:"-"`
	TrackingNumber       string           `json:"tracking_number,omitempty"`
	TrackingStatus       string           `json:"tracking_status,omitempty"`
	TrackingURLProvider  string           `json:"tracking_url_provider,omitempty"`
	Eta                  time.Time        `json:"eta,omitempty"`
	LabelURL             string           `json:"label_url,omitempty"`
	CommercialInvoiceURL string           `json:"commercial_invoice_url,omitempty"`
	Messages             []*OutputMessage `json:"messages,omitempty"`
	QRCodeURL            string           `json:"qr_code_url,omitempty"`
	Async                bool             `json:"async"`
}

// CoreRate is the simplified rate returned with an instant label purchase.
type CoreRate struct {
	ObjectID          string `json:"object_id"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	AmountLocal       string `json:"amount_local"`
	CurrencyLocal     string `json:"currency_local"`
	Provider          string `json:"provider"`
	CarrierAccount    string `json:"carrier_account"`
	ServiceLevelName  string `json:"servicelevel_name"`
	ServiceLevelToken string `json:"servicelevel_token"`
}

// UnmarshalJSON accepts both the rate ID and expanded rate response shapes.
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type transaction Transaction
	decoded := transaction(*t)
	wire := struct {
		*transaction
		Rate json.RawMessage `json:"rate"`
	}{transaction: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if len(wire.Rate) != 0 {
		decoded.RateDetails = nil
		decoded.Rate = ""
		if wire.Rate[0] == '{' {
			if err := json.Unmarshal(wire.Rate, &decoded.RateDetails); err != nil {
				return err
			}
			decoded.Rate = decoded.RateDetails.ObjectID
		} else if err := json.Unmarshal(wire.Rate, &decoded.Rate); err != nil {
			return err
		}
	}
	*t = Transaction(decoded)
	return nil
}

// TransactionListOptions filters transactions, including every parcel label
// purchased with the same multi-piece shipment rate.
type TransactionListOptions struct {
	Rate    string
	Results int // Results per page, up to 100; zero uses the client's default of 25.
}
