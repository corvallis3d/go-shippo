package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateShipment creates a new shipment object.
// Use Extra.IsReturn for return shipments, subject to the carrier's return terms.
func (c *Client) CreateShipment(input *models.ShipmentInput) (*models.Shipment, error) {
	return c.CreateShipmentContext(context.Background(), input)
}

// CreateShipmentContext is CreateShipment with cancellation and deadline support.
func (c *Client) CreateShipmentContext(ctx context.Context, input *models.ShipmentInput) (*models.Shipment, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if input.ShipmentDate.IsZero() {
		// Set the request default without mutating a caller's reusable input.
		copy := *input
		copy.ShipmentDate = time.Now()
		input = &copy
	}

	switch input.AddressFrom.(type) {
	case string, *models.AddressInput, nil:
	default:
		return nil, errors.New("invalid address_from data type")
	}

	switch input.AddressTo.(type) {
	case string, *models.AddressInput, nil:
	default:
		return nil, errors.New("invalid address_to data type")
	}

	switch input.AddressReturn.(type) {
	case string, *models.AddressInput, nil:
	default:
		return nil, errors.New("invalid address_return data type")
	}

	switch input.Parcels.(type) {
	case string, *models.ParcelInput, []string, []*models.ParcelInput, nil:
	default:
		return nil, errors.New("invalid parcels data type")
	}

	switch input.CustomsDeclaration.(type) {
	case string, *models.CustomsDeclarationInput, nil:
	default:
		return nil, errors.New("invalid customs_declaration data type")
	}

	output := &models.Shipment{}
	err := c.doContext(ctx, http.MethodPost, "/shipments/", input, output)
	return output, err
}

// RetrieveShipment retrieves an existing shipment by object id.
func (c *Client) RetrieveShipment(objectID string) (*models.Shipment, error) {
	return c.RetrieveShipmentContext(context.Background(), objectID)
}

// RetrieveShipmentContext is RetrieveShipment with cancellation and deadline support.
func (c *Client) RetrieveShipmentContext(ctx context.Context, objectID string) (*models.Shipment, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Shipment{}
	err := c.doContext(ctx, http.MethodGet, "/shipments/"+objectID, nil, output)
	return output, err
}

// ListAllShipments lists all shipment objects.
func (c *Client) ListAllShipments() ([]*models.Shipment, error) {
	return c.ListAllShipmentsContext(context.Background())
}

// ListAllShipmentsContext is ListAllShipments with cancellation and deadline support.
func (c *Client) ListAllShipmentsContext(ctx context.Context) ([]*models.Shipment, error) {
	list := []*models.Shipment{}
	err := c.doListContext(ctx, http.MethodGet, "/shipments/", nil, func(v json.RawMessage) error {
		item := &models.Shipment{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}
