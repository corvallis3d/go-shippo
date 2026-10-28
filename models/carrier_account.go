package models

// See https://docs.goshippo.com/api-reference/carrier-accounts/create-a-new-carrier-account
type CarrierAccountInput struct {
	Carrier    string                 `json:"carrier"`
	AccountID  string                 `json:"account_id"`
	Parameters map[string]interface{} `json:"parameters"`
	Active     bool                   `json:"active"`
}

// See https://docs.goshippo.com/api-reference/carrier-accounts/retrieve-a-carrier-account
type CarrierAccount struct {
	CarrierAccountInput
	CommonOutputFields
	Test            bool                          `json:"test"`
	CarrierName     string                        `json:"carrier_name,omitempty"`
	IsShippoAccount bool                          `json:"is_shippo_account"`
	Metadata        string                        `json:"metadata,omitempty"`
	ServiceLevels   []*CarrierAccountServiceLevel `json:"service_levels,omitempty"`
	ObjectInfo      *CarrierAccountObjectInfo     `json:"object_info,omitempty"`
}

type CarrierAccountServiceLevel struct {
	Name                 string `json:"name"`
	Token                string `json:"token"`
	SupportsReturnLabels bool   `json:"supports_return_labels"`
}

// CarrierAccountObjectInfo reports connection and verification state.
type CarrierAccountObjectInfo struct {
	Authentication *CarrierAccountAuthentication `json:"authentication,omitempty"`
	Verification   *CarrierAccountVerification   `json:"verification,omitempty"`
}

type CarrierAccountAuthentication struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type CarrierAccountVerification struct {
	Policy string `json:"policy"` // required or not-required
}
