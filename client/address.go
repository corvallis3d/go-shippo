package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateAddress creates a new address object.
func (c *Client) CreateAddress(input *models.AddressInput) (*models.Address, error) {
	return c.CreateAddressContext(context.Background(), input)
}

// CreateAddressContext is CreateAddress with cancellation and deadline support.
func (c *Client) CreateAddressContext(ctx context.Context, input *models.AddressInput) (*models.Address, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.Address{}
	err := c.doContext(ctx, http.MethodPost, "/addresses/", input, output)
	return output, err
}

// RetrieveAddress retrieves an existing address by object id.
func (c *Client) RetrieveAddress(objectID string) (*models.Address, error) {
	return c.RetrieveAddressContext(context.Background(), objectID)
}

// RetrieveAddressContext is RetrieveAddress with cancellation and deadline support.
func (c *Client) RetrieveAddressContext(ctx context.Context, objectID string) (*models.Address, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Address{}
	err := c.doContext(ctx, http.MethodGet, "/addresses/"+objectID, nil, output)
	return output, err
}

// ValidateAddress validates an existing address by object ID.
func (c *Client) ValidateAddress(objectID string) (*models.Address, error) {
	return c.ValidateAddressContext(context.Background(), objectID)
}

// ValidateAddressContext is ValidateAddress with cancellation and deadline support.
func (c *Client) ValidateAddressContext(ctx context.Context, objectID string) (*models.Address, error) {
	if objectID == "" {
		return nil, errors.New("empty object ID")
	}
	output := &models.Address{}
	err := c.doContext(ctx, http.MethodGet, "/addresses/"+objectID+"/validate", nil, output)
	return output, err
}

// ListAllAddresses lists all addresses.
func (c *Client) ListAllAddresses() ([]*models.Address, error) {
	return c.ListAllAddressesContext(context.Background())
}

// ListAllAddressesContext is ListAllAddresses with cancellation and deadline support.
func (c *Client) ListAllAddressesContext(ctx context.Context) ([]*models.Address, error) {
	list := []*models.Address{}
	err := c.doListContext(ctx, http.MethodGet, "/addresses/", nil, func(v json.RawMessage) error {
		item := &models.Address{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}
