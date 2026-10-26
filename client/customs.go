package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateCustomsItem creates a new customs item object.
func (c *Client) CreateCustomsItem(input *models.CustomsItemInput) (*models.CustomsItem, error) {
	return c.CreateCustomsItemContext(context.Background(), input)
}

// CreateCustomsItemContext is CreateCustomsItem with cancellation and deadline support.
func (c *Client) CreateCustomsItemContext(ctx context.Context, input *models.CustomsItemInput) (*models.CustomsItem, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.CustomsItem{}
	err := c.doContext(ctx, http.MethodPost, "/customs/items/", input, output)
	return output, err
}

// RetrieveCustomsItem retrieves an existing customs item by object id.
func (c *Client) RetrieveCustomsItem(objectID string) (*models.CustomsItem, error) {
	return c.RetrieveCustomsItemContext(context.Background(), objectID)
}

// RetrieveCustomsItemContext is RetrieveCustomsItem with cancellation and deadline support.
func (c *Client) RetrieveCustomsItemContext(ctx context.Context, objectID string) (*models.CustomsItem, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.CustomsItem{}
	err := c.doContext(ctx, http.MethodGet, "/customs/items/"+objectID, nil, output)
	return output, err
}

// ListAllCustomsItems lists all customs item objects.
func (c *Client) ListAllCustomsItems() ([]*models.CustomsItem, error) {
	return c.ListAllCustomsItemsContext(context.Background())
}

// ListAllCustomsItemsContext is ListAllCustomsItems with cancellation and deadline support.
func (c *Client) ListAllCustomsItemsContext(ctx context.Context) ([]*models.CustomsItem, error) {
	list := []*models.CustomsItem{}
	err := c.doListContext(ctx, http.MethodGet, "/customs/items/", nil, func(v json.RawMessage) error {
		item := &models.CustomsItem{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}

// CreateCustomsDeclaration creates a new customs declaration object.
func (c *Client) CreateCustomsDeclaration(input *models.CustomsDeclarationInput) (*models.CustomsDeclaration, error) {
	return c.CreateCustomsDeclarationContext(context.Background(), input)
}

// CreateCustomsDeclarationContext is CreateCustomsDeclaration with cancellation and deadline support.
func (c *Client) CreateCustomsDeclarationContext(ctx context.Context, input *models.CustomsDeclarationInput) (*models.CustomsDeclaration, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	switch input.Items.(type) {
	case []string, []*models.CustomsItemInput, nil:
	default:
		return nil, errors.New("invalid items data type")
	}

	output := &models.CustomsDeclaration{}
	err := c.doContext(ctx, http.MethodPost, "/customs/declarations/", input, output)
	return output, err
}

// RetrieveCustomsDeclaration retrieves an existing customs declaration by object id.
func (c *Client) RetrieveCustomsDeclaration(objectID string) (*models.CustomsDeclaration, error) {
	return c.RetrieveCustomsDeclarationContext(context.Background(), objectID)
}

// RetrieveCustomsDeclarationContext is RetrieveCustomsDeclaration with cancellation and deadline support.
func (c *Client) RetrieveCustomsDeclarationContext(ctx context.Context, objectID string) (*models.CustomsDeclaration, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.CustomsDeclaration{}
	err := c.doContext(ctx, http.MethodGet, "/customs/declarations/"+objectID, nil, output)
	return output, err
}

// ListAllCustomsDeclaration lists all customs declaration objects.
func (c *Client) ListAllCustomsDeclaration() ([]*models.CustomsDeclaration, error) {
	return c.ListAllCustomsDeclarationContext(context.Background())
}

// ListAllCustomsDeclarationContext is ListAllCustomsDeclaration with cancellation and deadline support.
func (c *Client) ListAllCustomsDeclarationContext(ctx context.Context) ([]*models.CustomsDeclaration, error) {
	list := []*models.CustomsDeclaration{}
	err := c.doListContext(ctx, http.MethodGet, "/customs/declarations/", nil, func(v json.RawMessage) error {
		item := &models.CustomsDeclaration{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}
