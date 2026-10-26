package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateCarrierAccount creates a new carrier account object.
func (c *Client) CreateCarrierAccount(input *models.CarrierAccountInput) (*models.CarrierAccount, error) {
	return c.CreateCarrierAccountContext(context.Background(), input)
}

// CreateCarrierAccountContext is CreateCarrierAccount with cancellation and deadline support.
func (c *Client) CreateCarrierAccountContext(ctx context.Context, input *models.CarrierAccountInput) (*models.CarrierAccount, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.CarrierAccount{}
	err := c.doContext(ctx, http.MethodPost, "/carrier_accounts/", input, output)
	return output, err
}

// RetrieveCarrierAccount retrieves an existing carrier account by object id.
func (c *Client) RetrieveCarrierAccount(objectID string) (*models.CarrierAccount, error) {
	return c.RetrieveCarrierAccountContext(context.Background(), objectID)
}

// RetrieveCarrierAccountContext is RetrieveCarrierAccount with cancellation and deadline support.
func (c *Client) RetrieveCarrierAccountContext(ctx context.Context, objectID string) (*models.CarrierAccount, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.CarrierAccount{}
	err := c.doContext(ctx, http.MethodGet, "/carrier_accounts/"+objectID, nil, output)
	return output, err
}

// ListAllCarrierAccounts lists all carrier accounts.
func (c *Client) ListAllCarrierAccounts() ([]*models.CarrierAccount, error) {
	return c.ListAllCarrierAccountsContext(context.Background())
}

// ListAllCarrierAccountsContext is ListAllCarrierAccounts with cancellation and deadline support.
func (c *Client) ListAllCarrierAccountsContext(ctx context.Context) ([]*models.CarrierAccount, error) {
	list := []*models.CarrierAccount{}
	err := c.doListContext(ctx, http.MethodGet, "/carrier_accounts/", nil, func(v json.RawMessage) error {
		item := &models.CarrierAccount{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}

// UpdateCarrierAccount updates an existing carrier account.
// AccountID and Carrier cannot be updated because they form the unique identifier together.
func (c *Client) UpdateCarrierAccount(objectID string, input *models.CarrierAccountInput) (*models.CarrierAccount, error) {
	return c.UpdateCarrierAccountContext(context.Background(), objectID, input)
}

// UpdateCarrierAccountContext is UpdateCarrierAccount with cancellation and deadline support.
func (c *Client) UpdateCarrierAccountContext(ctx context.Context, objectID string, input *models.CarrierAccountInput) (*models.CarrierAccount, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.CarrierAccount{}
	err := c.doContext(ctx, http.MethodPut, "/carrier_accounts/"+objectID, input, output)
	return output, err
}
