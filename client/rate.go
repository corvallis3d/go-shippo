package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// GetShippingRates gets rates for a shipping object.
// An empty currencyCode retrieves the existing rates without currency conversion.
func (c *Client) GetShippingRates(shipmentObjectID, currencyCode string) ([]*models.Rate, error) {
	return c.GetShippingRatesContext(context.Background(), shipmentObjectID, currencyCode)
}

// GetShippingRatesContext is GetShippingRates with cancellation and deadline support.
func (c *Client) GetShippingRatesContext(ctx context.Context, shipmentObjectID, currencyCode string) ([]*models.Rate, error) {
	if shipmentObjectID == "" {
		return nil, errors.New("Empty shipment object ID")
	}
	path := "/shipments/" + shipmentObjectID + "/rates/"
	if currencyCode != "" {
		path += currencyCode
	}

	list := []*models.Rate{}
	err := c.doListContext(ctx, http.MethodGet, path, nil, func(v json.RawMessage) error {
		item := &models.Rate{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}

// RetrieveRate retrieves an existing rate by object id.
func (c *Client) RetrieveRate(objectID string) (*models.Rate, error) {
	return c.RetrieveRateContext(context.Background(), objectID)
}

// RetrieveRateContext is RetrieveRate with cancellation and deadline support.
func (c *Client) RetrieveRateContext(ctx context.Context, objectID string) (*models.Rate, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Rate{}
	err := c.doContext(ctx, http.MethodGet, "/rates/"+objectID, nil, output)
	return output, err
}
