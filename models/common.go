package models

import (
	"encoding/json"
	"time"
)

const (
	CarrierAustraliaPost = "australia_post"
	CarrierCanadaPost    = "canada_post"
	CarrierDeutschePost  = "deutsche_post"
	CarrierDHLGermany    = "dhl_germany"
	CarrierDHLeCommerce  = "dhl_ecommerce"
	CarrierDHLExpress    = "dhl_express"
	CarrierFedEx         = "fedex"
	CarrierGLSGermany    = "gls_de"
	CarrierGLSFrance     = "gls_fr"
	CarrierLasership     = "lasership"
	CarrierMondialRelay  = "mondial_relay"
	CarrierNewgistics    = "newgistics"
	CarrierOnTrac        = "ontrac"
	CarrierPurolator     = "purolator"
	CarrierUPS           = "ups"
	CarrierUSPS          = "usps"

	ObjectStateValid   = "VALID"
	ObjectStateInvalid = "INVALID"

	StatusWaiting           = "WAITING"
	StatusQueued            = "QUEUED"
	StatusSuccess           = "SUCCESS"
	StatusError             = "ERROR"
	StatusRefunded          = "REFUNDED"
	StatusRefundPending     = "REFUNDPENDING"
	StatusRefundRejected    = "REFUNDREJECTED"
	StatusPending           = "PENDING"
	StatusValidating        = "VALIDATING"
	StatusValid             = "VALID"
	StatusInvalid           = "INVALID"
	StatusPurchasing        = "PURCHASING"
	StatusPurchased         = "PURCHASED"
	StatusIncomplete        = "INCOMPLETE"
	StatusTransactionFailed = "TRANSACTION_FAILED"

	InsuranceProviderFedEx  = "FEDEX"
	InsuranceProviderUPS    = "UPS"
	InsuranceProviderOnTrac = "ONTRAC"

	CODPaymentMethodSecuredFunds = "SECURED_FUNDS"
	CODPaymentMethodCash         = "CASH"
	CODPaymentMethodAny          = "ANY"

	DistanceUnitCentiMeter = "cm"
	DistanceUnitInch       = "in"
	DistanceUnitFoot       = "ft"
	DistanceUnitMilliMeter = "mm"
	DistanceUnitMeter      = "m"
	DistanceUnitYard       = "yd"

	MassUnitGram     = "g"
	MassUnitOunce    = "oz"
	MassUnitPound    = "lb"
	MassUnitKiloGram = "kg"

	LabelFileTypePNG      = "PNG"
	LabelFileTypePNG23X75 = "PNG_2.3x7.5"
	LabelFileTypePDF      = "PDF"
	LabelFileTypePDF23X75 = "PDF_2.3x7.5"
	LabelFileTypePDF4X6   = "PDF_4x6"
	LabelFileTypePDF4X8   = "PDF_4x8"
	LabelFileTypePDFA4    = "PDF_A4"
	LabelFileTypePDFA5    = "PDF_A5"
	LabelFileTypePDFA6    = "PDF_A6"
	LabelFileTypeZPLII    = "ZPLII"
)

type CommonOutputFields struct {
	ObjectCreated time.Time `json:"object_created,omitempty"`
	ObjectUpdated time.Time `json:"object_updated,omitempty"`
	ObjectID      string    `json:"object_id,omitempty"`
	ObjectOwner   string    `json:"object_owner,omitempty"`
}

type OutputMessage struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Text    string `json:"text,omitempty"`
	Source  string `json:"source,omitempty"`
}

type ListAPIOutput struct {
	Count           int               `json:"count"`
	NextPageURL     *string           `json:"next"`
	PreviousPageURL *string           `json:"previous"`
	Results         []json.RawMessage `json:"results"`
}
