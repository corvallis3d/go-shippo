package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateOrder creates a new order with specified details.
func (c *Client) CreateOrder(input *models.OrderInput) (*models.OrderResponse, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}
	output := &models.OrderResponse{}
	err := c.do(http.MethodPost, "/orders/", input, output)
	return output, err
}

// GetOrder retrieves an existing order by its ID.
func (c *Client) GetOrder(objectID string) (*models.OrderResponse, error) {
	if objectID == "" {
		return nil, errors.New("empty order ID")
	}
	output := &models.OrderResponse{}
	err := c.do(http.MethodGet, fmt.Sprintf("/orders/%s", objectID), nil, output)
	return output, err
}

// ListOrders lists all orders with optional filters.
// Filters for shop app, date range, and order status are supported.
func (c *Client) ListOrders(options *models.OrderListOptions) (*models.OrderListResponse, error) {
	url := "/orders/"
	queryParams := []string{}
	if options != nil {
		if options.ShopApp != "" {
			queryParams = append(queryParams, "shop_app="+options.ShopApp)
		}
		if !options.StartDate.IsZero() {
			queryParams = append(queryParams, "start_date="+options.StartDate.Format(time.RFC3339))
		}
		if !options.EndDate.IsZero() {
			queryParams = append(queryParams, "end_date="+options.EndDate.Format(time.RFC3339))
		}
		for _, status := range options.OrderStatus {
			queryParams = append(queryParams, "order_status[]="+status)
		}
		if options.Page > 0 {
			queryParams = append(queryParams, fmt.Sprintf("page=%d", options.Page))
		}
		if options.Results > 0 {
			queryParams = append(queryParams, fmt.Sprintf("results=%d", options.Results))
		}
	}
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	output := &models.OrderListResponse{}
	err := c.do(http.MethodGet, url, nil, output)
	return output, err
}
