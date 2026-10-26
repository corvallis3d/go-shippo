package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateOrder creates a new order with specified details.
func (c *Client) CreateOrder(input *models.OrderInput) (*models.OrderResponse, error) {
	return c.CreateOrderContext(context.Background(), input)
}

// CreateOrderContext is CreateOrder with cancellation and deadline support.
func (c *Client) CreateOrderContext(ctx context.Context, input *models.OrderInput) (*models.OrderResponse, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}
	output := &models.OrderResponse{}
	err := c.doContext(ctx, http.MethodPost, "/orders/", input, output)
	return output, err
}

// GetOrder retrieves an existing order by its ID.
func (c *Client) GetOrder(objectID string) (*models.OrderResponse, error) {
	return c.GetOrderContext(context.Background(), objectID)
}

// GetOrderContext is GetOrder with cancellation and deadline support.
func (c *Client) GetOrderContext(ctx context.Context, objectID string) (*models.OrderResponse, error) {
	if objectID == "" {
		return nil, errors.New("empty order ID")
	}
	output := &models.OrderResponse{}
	err := c.doContext(ctx, http.MethodGet, fmt.Sprintf("/orders/%s", objectID), nil, output)
	return output, err
}

// ListOrders retrieves one page of orders with optional filters.
// Filters for shop app, date range, and order status are supported.
func (c *Client) ListOrders(options *models.OrderListOptions) (*models.OrderListResponse, error) {
	return c.ListOrdersContext(context.Background(), options)
}

// ListOrdersContext is ListOrders with cancellation and deadline support.
func (c *Client) ListOrdersContext(ctx context.Context, options *models.OrderListOptions) (*models.OrderListResponse, error) {
	path := "/orders/"
	queryParams := url.Values{}
	if options != nil {
		if options.ShopApp != "" {
			queryParams.Set("shop_app", options.ShopApp)
		}
		if !options.StartDate.IsZero() {
			queryParams.Set("start_date", options.StartDate.UTC().Format(time.RFC3339))
		}
		if !options.EndDate.IsZero() {
			queryParams.Set("end_date", options.EndDate.UTC().Format(time.RFC3339))
		}
		for _, status := range options.OrderStatus {
			queryParams.Add("order_status[]", status)
		}
		if options.Page > 0 {
			queryParams.Set("page", strconv.Itoa(options.Page))
		}
		if options.Results > 0 {
			queryParams.Set("results", strconv.Itoa(options.Results))
		}
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	output := &models.OrderListResponse{}
	err := c.doContext(ctx, http.MethodGet, path, nil, output)
	return output, err
}
