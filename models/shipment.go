package models

import (
	"encoding/json"
	"time"
)

const (
	ShipmentExtraSignatureConfirmationStandard  = "STANDARD"
	ShipmentExtraSignatureConfirmationAdult     = "ADULT"
	ShipmentExtraSignatureConfirmationCertified = "CERTIFIED"
	ShipmentExtraSignatureConfirmationIndirect  = "INDIRECT"
	ShipmentExtraSignatureConfirmationCarrier   = "CARRIER_CONFIRMATION"

	ShipmentBillingTypeSender     = "SENDER"
	ShipmentBillingTypeRecipient  = "RECIPIENT"
	ShipmentBillingTypeThirdParty = "THIRD_PARTY"

	ShipmentUSPSSortTypeNDC          = "NDC"
	ShipmentUSPSSortTypeFiveDigit    = "FiveDigit"
	ShipmentUSPSSortTypeMixedNDC     = "MixedNDC"
	ShipmentUSPSSortTypeNonPresorted = "Nonpresorted"
	ShipmentUSPSSortTypePresorted    = "Presorted"
	ShipmentUSPSSortTypeSCF          = "SCF"
	ShipmentUSPSSortTypeSinglePiece  = "SinglePiece"
	ShipmentUSPSSortTypeThreeDigit   = "ThreeDigit"

	ShipmentUSPSEntryFacilityDNDC  = "DNDC"
	ShipmentUSPSEntryFacilityDDU   = "DDU"
	ShipmentUSPSEntryFacilityDSCF  = "DSCF"
	ShipmentUSPSEntryFacilityONDC  = "ONDC"
	ShipmentUSPSEntryFacilityOther = "Other"

	ShipmentDangerousGoodsCode01 = "01"
	ShipmentDangerousGoodsCode02 = "02"
	ShipmentDangerousGoodsCode03 = "03"
	ShipmentDangerousGoodsCode04 = "04"
	ShipmentDangerousGoodsCode05 = "05"
	ShipmentDangerousGoodsCode06 = "06"
	ShipmentDangerousGoodsCode07 = "07"
	ShipmentDangerousGoodsCode08 = "08"
	ShipmentDangerousGoodsCode09 = "09"
)

// See https://docs.goshippo.com/api-reference/shipments/create-a-new-shipment
type ShipmentInput struct {
	AddressFrom        interface{}    `json:"address_from"`                  // string or *AddressInput
	AddressTo          interface{}    `json:"address_to"`                    // string or *AddressInput
	Parcels            interface{}    `json:"parcels"`                       // string, []string, *ParcelInput, or []*ParcelInput
	AddressReturn      interface{}    `json:"address_return,omitempty"`      // string or *AddressInput
	CustomsDeclaration interface{}    `json:"customs_declaration,omitempty"` // string or *CustomsDeclarationInput
	ShipmentDate       time.Time      `json:"shipment_date,omitempty"`
	CarrierAccounts    []string       `json:"carrier_accounts,omitempty"`
	Metadata           string         `json:"metadata,omitempty"`
	Extra              *ShipmentExtra `json:"extra,omitempty"`
	Async              bool           `json:"async"`
}

type ShipmentExtra struct {
	QRCodeRequested                   bool                        `json:"qr_code_requested,omitempty"`
	SignatureConfirmation             string                      `json:"signature_confirmation,omitempty"`
	SaturdayDelivery                  bool                        `json:"saturday_delivery"`
	BypassAddressValidation           bool                        `json:"bypass_address_validation"`
	RequestRetailRates                bool                        `json:"request_retail_rates"`
	CustomerBranch                    string                      `json:"customer_branch,omitempty"`
	Billing                           *ShipmentBilling            `json:"billing,omitempty"`
	COD                               *ShipmentCOD                `json:"COD,omitempty"`
	Insurance                         *ShipmentInsurance          `json:"insurance,omitempty"`
	USPSSortType                      string                      `json:"usps_sort_type,omitempty"`
	USPSEntryFacility                 string                      `json:"usps_entry_facility,omitempty"`
	DangerousGoodsCode                string                      `json:"dangerous_goods_code,omitempty"`
	IsReturn                          bool                        `json:"is_return,omitempty"`
	Reference1                        string                      `json:"reference_1,omitempty"`
	Reference2                        string                      `json:"reference_2,omitempty"`
	Reference3                        string                      `json:"reference_3,omitempty"`
	Reference4                        string                      `json:"reference_4,omitempty"`
	HealthcareHandling                *ShipmentHealthcareHandling `json:"healthcare_handling,omitempty"`
	CustomerReference                 *ShipmentLabelReference     `json:"customer_reference,omitempty"`
	PONumber                          *ShipmentLabelReference     `json:"po_number,omitempty"`
	InvoiceNumber                     *ShipmentLabelReference     `json:"invoice_number,omitempty"`
	DeptNumber                        *ShipmentLabelReference     `json:"dept_number,omitempty"`
	RMANumber                         *ShipmentLabelReference     `json:"rma_number,omitempty"`
	AccountsReceivableCustomerAccount *ShipmentLabelReference     `json:"accounts_receivable_customer_account,omitempty"`
	AppropriationNumber               *ShipmentLabelReference     `json:"appropriation_number,omitempty"`
	BillOfLadingNumber                *ShipmentLabelReference     `json:"bill_of_lading_number,omitempty"`
	CODNumber                         *ShipmentLabelReference     `json:"cod_number,omitempty"`
	DealerOrderNumber                 *ShipmentLabelReference     `json:"dealer_order_number,omitempty"`
	FDAProductCode                    *ShipmentLabelReference     `json:"fda_product_code,omitempty"`
	ManifestNumber                    *ShipmentLabelReference     `json:"manifest_number,omitempty"`
	ModelNumber                       *ShipmentLabelReference     `json:"model_number,omitempty"`
	PartNumber                        *ShipmentLabelReference     `json:"part_number,omitempty"`
	ProductionCode                    *ShipmentLabelReference     `json:"production_code,omitempty"`
	PurchaseRequestNumber             *ShipmentLabelReference     `json:"purchase_request_number,omitempty"`
	SalespersonNumber                 *ShipmentLabelReference     `json:"salesperson_number,omitempty"`
	SerialNumber                      *ShipmentLabelReference     `json:"serial_number,omitempty"`
	StoreNumber                       *ShipmentLabelReference     `json:"store_number,omitempty"`
	TransactionReferenceNumber        *ShipmentLabelReference     `json:"transaction_reference_number,omitempty"`
	DangerousGoods                    *ShipmentDangerousGoods     `json:"dangerous_goods,omitempty"`
	Alcohol                           *ShipmentAlcohol            `json:"alcohol,omitempty"`
	DryIce                            *ShipmentDryIce             `json:"dry_ice,omitempty"`
	AncillaryEndorsement              string                      `json:"ancillary_endorsement,omitempty"`
	AuthorityToLeave                  *bool                       `json:"authority_to_leave,omitempty"`
	CarbonNeutral                     bool                        `json:"carbon_neutral,omitempty"`
	CarrierHubID                      string                      `json:"carrier_hub_id,omitempty"`
	CarrierHubTravelTime              int                         `json:"carrier_hub_travel_time,omitempty"`
	ContainerType                     string                      `json:"container_type,omitempty"`
	CriticalPullTime                  string                      `json:"critical_pull_time,omitempty"`
	DeliveryInstructions              string                      `json:"delivery_instructions,omitempty"`
	FulfillmentCenter                 string                      `json:"fulfillment_center,omitempty"`
	LasershipAttrs                    []string                    `json:"lasership_attrs,omitempty"`
	LasershipDeclaredValue            string                      `json:"lasership_declared_value,omitempty"`
	PreferredDeliveryTimeframe        string                      `json:"preferred_delivery_timeframe,omitempty"`
	Premium                           bool                        `json:"premium,omitempty"`
	ReturnServiceType                 string                      `json:"return_service_type,omitempty"`
	USMCAEligible                     bool                        `json:"usmca_eligible,omitempty"`
}

// ShipmentHealthcareHandling requests a carrier's healthcare handling program.
// The carrier determines account and service eligibility.
type ShipmentHealthcareHandling struct {
	Enabled bool `json:"enabled"`
}

// ShipmentLabelReference is the reference object used by FedEx and UPS label
// fields. A nil Prefix retains the carrier default; an empty prefix removes it.
// Carrier-specific length, format and sorting rules are enforced by Shippo.
type ShipmentLabelReference struct {
	Prefix  *string `json:"prefix,omitempty"`
	Value   string  `json:"value,omitempty"`
	RefSort *int    `json:"ref_sort,omitempty"`
}

type ShipmentDangerousGoods struct {
	Contains           bool                            `json:"contains"`
	BiologicalMaterial *ShipmentDangerousGoodsContents `json:"biological_material,omitempty"`
	LithiumBatteries   *ShipmentDangerousGoodsContents `json:"lithium_batteries,omitempty"`
}

type ShipmentDangerousGoodsContents struct {
	Contains bool `json:"contains"`
}

type ShipmentAlcohol struct {
	ContainsAlcohol bool   `json:"contains_alcohol"`
	RecipientType   string `json:"recipient_type,omitempty"`
}

type ShipmentDryIce struct {
	ContainsDryIce bool   `json:"contains_dry_ice"`
	Weight         string `json:"weight"` // Kilograms, as specified by Shippo.
}

type ShipmentCOD struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

type ShipmentInsurance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Provider string `json:"provider,omitempty"` // Omit for Shippo's third-party insurance.
	Content  string `json:"content"`
}

type ShipmentBilling struct {
	Type    string `json:"type"`
	Account string `json:"account"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

// See https://docs.goshippo.com/api-reference/shipments/retrieve-a-shipment
type Shipment struct {
	CommonOutputFields
	Metadata           string              `json:"metadata,omitempty"`
	Extra              *ShipmentExtra      `json:"extra,omitempty"`
	CarrierAccounts    []string            `json:"carrier_accounts,omitempty"`
	ShipmentDate       time.Time           `json:"shipment_date,omitempty"`
	Status             string              `json:"status,omitempty"`
	AddressFrom        *Address            `json:"address_from,omitempty"`
	AddressTo          *Address            `json:"address_to,omitempty"`
	Parcels            []*Parcel           `json:"parcels,omitempty"`
	AddressReturn      *Address            `json:"address_return,omitempty"`
	CustomsDeclaration *CustomsDeclaration `json:"customs_declaration,omitempty"`
	Rates              []*Rate             `json:"rates,omitempty"`
	Messages           []*OutputMessage    `json:"messages,omitempty"`
	Test               bool                `json:"test,omitempty"`
}

// MarshalJSON normalizes single parcels to arrays and omits an unset shipment
// date. time.Time's zero value is not omitted by encoding/json's omitempty.
func (s ShipmentInput) MarshalJSON() ([]byte, error) {
	type shipmentInput ShipmentInput
	switch parcels := s.Parcels.(type) {
	case string:
		s.Parcels = []string{parcels}
	case *ParcelInput:
		if parcels != nil {
			s.Parcels = []*ParcelInput{parcels}
		}
	}
	wire := struct {
		shipmentInput
		ShipmentDate *time.Time `json:"shipment_date,omitempty"`
	}{shipmentInput: shipmentInput(s)}
	if !s.ShipmentDate.IsZero() {
		wire.ShipmentDate = &s.ShipmentDate
	}
	return json.Marshal(wire)
}

// UnmarshalJSON accepts either a customs declaration ID or an expanded object.
func (s *Shipment) UnmarshalJSON(data []byte) error {
	type shipment Shipment
	decoded := shipment(*s)
	wire := struct {
		*shipment
		CustomsDeclaration json.RawMessage `json:"customs_declaration"`
	}{shipment: &decoded}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if len(wire.CustomsDeclaration) != 0 {
		decoded.CustomsDeclaration = nil
		if wire.CustomsDeclaration[0] == '"' {
			var id string
			if err := json.Unmarshal(wire.CustomsDeclaration, &id); err != nil {
				return err
			}
			decoded.CustomsDeclaration = &CustomsDeclaration{CommonOutputFields: CommonOutputFields{ObjectID: id}}
		} else if err := json.Unmarshal(wire.CustomsDeclaration, &decoded.CustomsDeclaration); err != nil {
			return err
		}
	}
	*s = Shipment(decoded)
	return nil
}
