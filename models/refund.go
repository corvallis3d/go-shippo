package models

// See https://docs.goshippo.com/api-reference/refunds/create-a-refund
type RefundInput struct {
	Transaction string `json:"transaction"`
	Async       bool   `json:"async"`
}

// See https://docs.goshippo.com/api-reference/refunds/retrieve-a-refund
type Refund struct {
	RefundInput
	CommonOutputFields
	Status string `json:"status,omitempty"`
	Test   bool   `json:"test"`
}
