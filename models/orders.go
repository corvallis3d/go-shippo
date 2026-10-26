package models

import (
	"encoding/json"
	"time"
)

// OrderInput represents the structure of an order creation request
type OrderInput struct {
	ToAddress            Address    `json:"to_address"`             // Required
	FromAddress          *Address   `json:"from_address,omitempty"` // Optional
	LineItems            []LineItem `json:"line_items,omitempty"`   // Optional
	PlacedAt             time.Time  `json:"placed_at"`              // Required
	OrderNumber          string     `json:"order_number,omitempty"`
	OrderStatus          string     `json:"order_status,omitempty"`
	ShippingCost         string     `json:"shipping_cost,omitempty"`
	ShippingCostCurrency string     `json:"shipping_cost_currency,omitempty"`
	ShippingMethod       string     `json:"shipping_method,omitempty"`
	SubtotalPrice        string     `json:"subtotal_price,omitempty"`
	TotalPrice           string     `json:"total_price,omitempty"`
	TotalTax             string     `json:"total_tax,omitempty"`
	Currency             string     `json:"currency,omitempty"` // Required if `total_price` is provided
	Weight               string     `json:"weight,omitempty"`
	WeightUnit           string     `json:"weight_unit,omitempty"`
	Notes                string     `json:"notes,omitempty"`
}

// OrderResponse represents the structure of an individual order in response
type OrderResponse struct {
	ObjectID             string        `json:"object_id"`
	ObjectOwner          string        `json:"object_owner"`
	OrderNumber          string        `json:"order_number"`
	OrderStatus          string        `json:"order_status"`
	PlacedAt             time.Time     `json:"placed_at"`
	ToAddress            Address       `json:"to_address"`
	FromAddress          *Address      `json:"from_address,omitempty"`
	LineItems            []LineItem    `json:"line_items"`
	ShippingCost         string        `json:"shipping_cost,omitempty"`
	ShippingCostCurrency string        `json:"shipping_cost_currency,omitempty"`
	ShippingMethod       string        `json:"shipping_method,omitempty"`
	SubtotalPrice        string        `json:"subtotal_price,omitempty"`
	TotalPrice           string        `json:"total_price"`
	TotalTax             string        `json:"total_tax,omitempty"`
	Currency             string        `json:"currency"`
	Weight               string        `json:"weight,omitempty"`
	WeightUnit           string        `json:"weight_unit,omitempty"`
	Notes                string        `json:"notes,omitempty"`
	ShopApp              string        `json:"shop_app,omitempty"`
	Transactions         []Transaction `json:"transactions,omitempty"`
	// These response fields are documented in Shippo's API changelog. The
	// current schema does not define the error payload's structure.
	ValidationStatus string          `json:"validation_status,omitempty"`
	Error            json.RawMessage `json:"error,omitempty"`
}

// OrderListResponse represents the structure of a paginated response for listing orders
type OrderListResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next,omitempty"`     // URL to the next page, if available
	Previous string          `json:"previous,omitempty"` // URL to the previous page, if available
	Results  []OrderResponse `json:"results"`            // Array of orders in the current page
}

// MarshalJSON keeps the existing Address-based Go API while sending only the
// writable address fields expected by order creation.
func (o OrderInput) MarshalJSON() ([]byte, error) {
	type orderInput OrderInput
	wire := struct {
		orderInput
		ToAddress   AddressInput  `json:"to_address"`
		FromAddress *AddressInput `json:"from_address,omitempty"`
		LineItems   []LineItem    `json:"line_items,omitempty"`
	}{orderInput: orderInput(o), ToAddress: o.ToAddress.AddressInput}
	if o.FromAddress != nil {
		wire.FromAddress = &o.FromAddress.AddressInput
	}
	if len(o.LineItems) != 0 {
		wire.LineItems = append([]LineItem(nil), o.LineItems...)
		for i := range wire.LineItems {
			// IDs belong to returned line items, not the create request.
			wire.LineItems[i].ObjectID = ""
		}
	}
	return json.Marshal(wire)
}

// LineItem represents an item in an order
type LineItem struct {
	ObjectID           string     `json:"object_id,omitempty"`
	Title              string     `json:"title,omitempty"`
	VariantTitle       string     `json:"variant_title,omitempty"`
	SKU                string     `json:"sku,omitempty"`
	Quantity           int        `json:"quantity,omitempty"`
	TotalPrice         string     `json:"total_price,omitempty"`
	Currency           string     `json:"currency,omitempty"`
	Weight             string     `json:"weight,omitempty"`
	WeightUnit         string     `json:"weight_unit,omitempty"`
	ManufactureCountry string     `json:"manufacture_country,omitempty"`
	MaxShipTime        *time.Time `json:"max_ship_time,omitempty"`
	MaxDeliveryTime    *time.Time `json:"max_delivery_time,omitempty"`
	Description        *string    `json:"description,omitempty"`
}

// OrderListOptions represents filtering options for listing orders
type OrderListOptions struct {
	Page        int       `json:"page,omitempty"`         // The page number to retrieve (pagination)
	Results     int       `json:"results,omitempty"`      // The number of results per page, max 100
	OrderStatus []string  `json:"order_status,omitempty"` // Array of order statuses to filter by
	ShopApp     string    `json:"shop_app,omitempty"`     // Filter by the originating platform (e.g., "Shopify", "WooCommerce")
	StartDate   time.Time `json:"start_date,omitempty"`   // Filter orders placed after this date
	EndDate     time.Time `json:"end_date,omitempty"`     // Filter orders placed before this date
}
