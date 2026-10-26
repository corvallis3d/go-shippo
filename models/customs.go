package models

const (
	CustomsNonDeliveryOptionAbandon = "ABANDON"
	CustomsNonDeliveryOptionReturn  = "RETURN"

	CustomsContentsTypeDocuments            = "DOCUMENTS"
	CustomsContentsTypeMerchandise          = "MERCHANDISE"
	CustomsContentsTypeGift                 = "GIFT"
	CustomsContentsTypeSample               = "SAMPLE"
	CustomsContentsTypeHumanitarianDonation = "HUMANITARIAN_DONATION"
	CustomsContentsTypeReturnMerchandise    = "RETURN_MERCHANDISE"
	CustomsContentsTypeOther                = "OTHER"

	CustomsEEL_PFC_NOEEI_30_37_A = "NOEEI_30_37_a"
	CustomsEEL_PFC_NOEEI_30_37_H = "NOEEI_30_37_h"
	CustomsEEL_PFC_NOEEI_30_37_F = "NOEEI_30_37_f"
	CustomsEEL_PFC_NOEEI_30_36   = "NOEEI_30_36"
	CustomsEEL_PFC_AES_ITN       = "AES_ITN"

	CustomsIncotermDDP = "DDP"
	CustomsIncotermDDU = "DDU"
	CustomsIncotermCPT = "CPT"
	CustomsIncotermCIP = "CIP"

	CustomsB13AFiledElectronically = "FILED_ELECTRONICALLY"
	CustomsB13ASummaryReporting    = "SUMMARY_REPORTING"
	CustomsB13ANotRequired         = "NOT_REQUIRED"

	CustomsTaxIDTypeEIN  = "EIN"
	CustomsTaxIDTypeVAT  = "VAT"
	CustomsTaxIDTypeIOSS = "IOSS"
	CustomsTaxIDTypeARN  = "ARN"
)

// See https://goshippo.com/docs/reference#customsitems
type CustomsItemInput struct {
	Description              string `json:"description"`
	Quantity                 int    `json:"quantity"`
	NetWeight                string `json:"net_weight"`
	MassUnit                 string `json:"mass_unit"`
	ValueAmount              string `json:"value_amount"`
	ValueCurrency            string `json:"value_currency"`
	OriginCountry            string `json:"origin_country"`
	TariffNumber             string `json:"tariff_number,omitempty"`
	HSCode                   string `json:"hs_code,omitempty"`
	SKUCode                  string `json:"sku_code,omitempty"`
	ECCNEAR99                string `json:"eccn_ear99,omitempty"`
	ManufacturerCode         string `json:"manufacturer_code,omitempty"`
	ManufacturerStandardCode string `json:"manufacturer_standard_code,omitempty"`
	EUExemptCategory         bool   `json:"eu_exempt_category,omitempty"`
	Metadata                 string `json:"metadata,omitempty"`
}

// See https://goshippo.com/docs/reference#customsitems
type CustomsItem struct {
	CustomsItemInput
	CommonOutputFields
	ObjectState string `json:"object_state,omitempty"`
	Test        bool   `json:"test"`
}

// See https://goshippo.com/docs/reference#customsdeclarations
type CustomsDeclarationInput struct {
	CertifySigner          string                         `json:"certify_signer"`
	Certify                bool                           `json:"certify"`
	Items                  interface{}                    `json:"items"` // []string or []*CustomsItemInput
	NonDeliveryOption      string                         `json:"non_delivery_option"`
	ContentsType           string                         `json:"contents_type"`
	ContentsExplanation    string                         `json:"contents_explanation,omitempty"`
	ExporterReference      string                         `json:"exporter_reference,omitempty"`
	ImporterReference      string                         `json:"importer_reference,omitempty"`
	Invoice                string                         `json:"invoice,omitempty"`
	License                string                         `json:"license,omitempty"`
	Certificate            string                         `json:"certificate,omitempty"`
	Notes                  string                         `json:"notes,omitempty"`
	EEL_PFC                string                         `json:"eel_pfc,omitempty"`
	AES_ITN                string                         `json:"aes_itn,omitempty"`
	Incoterm               string                         `json:"incoterm,omitempty"`
	Metadata               string                         `json:"metadata,omitempty"`
	Disclaimer             string                         `json:"disclaimer,omitempty"`
	RecipientIsBusiness    bool                           `json:"recipient_is_business,omitempty"`
	B13AFilingOption       string                         `json:"b13a_filing_option,omitempty"`
	B13ANumber             string                         `json:"b13a_number,omitempty"`
	CommercialInvoice      bool                           `json:"commercial_invoice,omitempty"`
	DutiesPayor            *CustomsDutiesPayor            `json:"duties_payor,omitempty"`
	ExporterIdentification *CustomsExporterIdentification `json:"exporter_identification,omitempty"`
	IsVATCollected         bool                           `json:"is_vat_collected,omitempty"`
}

type CustomsDutiesPayor struct {
	Account string                     `json:"account,omitempty"`
	Type    string                     `json:"type,omitempty"`
	Address *CustomsDutiesPayorAddress `json:"address,omitempty"`
}

type CustomsDutiesPayorAddress struct {
	Name    string `json:"name,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

type CustomsExporterIdentification struct {
	EORINumber string                    `json:"eori_number,omitempty"`
	TaxID      *CustomsTaxIdentification `json:"tax_id,omitempty"`
}

type CustomsTaxIdentification struct {
	Number string `json:"number,omitempty"`
	Type   string `json:"type,omitempty"`
}

// See https://goshippo.com/docs/reference#customsdeclarations
type CustomsDeclaration struct {
	CustomsDeclarationInput
	CommonOutputFields
	ObjectState string `json:"object_state,omitempty"`
	Test        bool   `json:"test"`
}
