package models

import "time"

type OrderInput struct {
	ToAddress            interface{} `json:"to_address"`
	LineItems            []LineItem  `json:"line_items"`
	PlacedAt             time.Time   `json:"placed_at"`
	OrderNumber          string      `json:"order_number"`
	OrderStatus          string      `json:"order_status"`
	ShippingCost         string      `json:"shipping_cost"`
	ShippingCostCurrency string      `json:"shipping_cost_currency"`
	ShippingMethod       string      `json:"shipping_method"`
	SubtotalPrice        string      `json:"subtotal_price"`
	TotalPrice           string      `json:"total_price"`
	TotalTax             string      `json:"total_tax"`
	Currency             string      `json:"currency"`
	Weight               string      `json:"weight"`
	WeightUnit           string      `json:"weight_unit"`
}

type LineItem struct {
	ObjectID           string    `json:"object_id"`
	Title              string    `json:"title"`
	VariantTitle       string    `json:"variant_title"`
	Sku                string    `json:"sku"`
	Quantity           int       `json:"quantity"`
	TotalPrice         string    `json:"total_price"`
	Currency           string    `json:"currency"`
	Weight             string    `json:"weight"`
	WeightUnit         string    `json:"weight_unit"`
	ManufactureCountry string    `json:"manufacture_country"`
	MaxShipTime        time.Time `json:"max_ship_time"`
	MaxDeliveryTime    time.Time `json:"max_delivery_time"`
}

// OrderListResponse represents the structure of a paginated response for listing orders
type OrderListResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next,omitempty"`
	Previous string          `json:"previous,omitempty"`
	Results  []OrderResponse `json:"results"`
}

type OrderResponse struct {
	ObjectID             string        `json:"object_id"`
	ObjectOwner          string        `json:"object_owner"`
	OrderNumber          string        `json:"order_number"`
	OrderStatus          string        `json:"order_status"`
	PlacedAt             time.Time     `json:"placed_at"`
	ToAddress            Address       `json:"to_address"`
	FromAddress          Address       `json:"from_address"` // Given it's null in the sample, this can be another struct or a pointer to a struct
	LineItems            []LineItem    `json:"line_items"`
	Items                []interface{} `json:"items"`
	Hidden               bool          `json:"hidden"`
	ShippingCost         string        `json:"shipping_cost"`
	ShippingCostCurrency string        `json:"shipping_cost_currency"`
	ShippingMethod       string        `json:"shipping_method"`
	ShopApp              string        `json:"shop_app"`
	SubtotalPrice        string        `json:"subtotal_price"`
	TotalPrice           string        `json:"total_price"`
	TotalTax             string        `json:"total_tax"`
	Currency             string        `json:"currency"`
	Transactions         []Transaction `json:"transactions"`
	Weight               string        `json:"weight"`
	WeightUnit           string        `json:"weight_unit"`
	Notes                interface{}   `json:"notes"` // Could be another struct or type based on actual data
}

// OrderListOptions represents filtering options for listing orders
type OrderListOptions struct {
	ShopApp     string    `json:"shop_app,omitempty"`
	StartDate   time.Time `json:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty"`
	OrderStatus []string  `json:"order_status,omitempty"`
	Page        int       `json:"page,omitempty"`
	Results     int       `json:"results,omitempty"`
}

// ShipmentRequest represents the structure for purchasing a label for an order
type ShipmentRequest struct {
	Shipment struct {
		AddressFrom Address  `json:"address_from"`
		AddressTo   string   `json:"address_to"`
		Parcels     []Parcel `json:"parcels"`
	} `json:"shipment"`
	CarrierAccount    string `json:"carrier_account"`
	ServiceLevelToken string `json:"servicelevel_token"`
	Order             string `json:"order"`
}

// TransactionResponse represents the response structure for a label purchase
type TransactionResponse struct {
	ObjectID       string `json:"object_id"`
	LabelURL       string `json:"label_url"`
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"`
}

// PackingSlipResponse represents the structure of a packing slip response
type PackingSlipResponse struct {
	SlipURL string    `json:"slip_url"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}
