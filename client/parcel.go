package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// CreateParcel creates a new parcel object.
func (c *Client) CreateParcel(input *models.ParcelInput) (*models.Parcel, error) {
	return c.CreateParcelContext(context.Background(), input)
}

// CreateParcelContext is CreateParcel with cancellation and deadline support.
func (c *Client) CreateParcelContext(ctx context.Context, input *models.ParcelInput) (*models.Parcel, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.Parcel{}
	err := c.doContext(ctx, http.MethodPost, "/parcels/", input, output)
	return output, err
}

// RetrieveParcel retrieves an existing parcel by object id.
func (c *Client) RetrieveParcel(objectID string) (*models.Parcel, error) {
	return c.RetrieveParcelContext(context.Background(), objectID)
}

// RetrieveParcelContext is RetrieveParcel with cancellation and deadline support.
func (c *Client) RetrieveParcelContext(ctx context.Context, objectID string) (*models.Parcel, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Parcel{}
	err := c.doContext(ctx, http.MethodGet, "/parcels/"+objectID, nil, output)
	return output, err
}

// ListAllParcels lists all parcel objects.
func (c *Client) ListAllParcels() ([]*models.Parcel, error) {
	return c.ListAllParcelsContext(context.Background())
}

// ListAllParcelsContext is ListAllParcels with cancellation and deadline support.
func (c *Client) ListAllParcelsContext(ctx context.Context) ([]*models.Parcel, error) {
	list := []*models.Parcel{}
	err := c.doListContext(ctx, http.MethodGet, "/parcels/", nil, func(v json.RawMessage) error {
		item := &models.Parcel{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}
