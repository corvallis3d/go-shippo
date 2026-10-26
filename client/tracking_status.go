package client

import (
	"context"
	"errors"
	"net/http"

	"github.com/corvallis3d/go-shippo/models"
)

// GetTrackingUpdate requests the tracking status of a shipment.
func (c *Client) GetTrackingUpdate(carrier, trackingNumber string) (*models.TrackingStatus, error) {
	return c.GetTrackingUpdateContext(context.Background(), carrier, trackingNumber)
}

// GetTrackingUpdateContext is GetTrackingUpdate with cancellation and deadline support.
func (c *Client) GetTrackingUpdateContext(ctx context.Context, carrier, trackingNumber string) (*models.TrackingStatus, error) {
	if carrier == "" {
		return nil, errors.New("Empty carrier")
	}
	if trackingNumber == "" {
		return nil, errors.New("Empty tracking number")
	}

	output := &models.TrackingStatus{}
	err := c.doContext(ctx, http.MethodGet, "/tracks/"+carrier+"/"+trackingNumber, nil, output)
	return output, err
}

// RegisterTrackingWebhook registers a carrier/tracking-number pair for tracking
// updates sent to the account's configured tracking webhooks.
func (c *Client) RegisterTrackingWebhook(carrier, trackingNumber, metadata string) (*models.TrackingStatus, error) {
	return c.RegisterTrackingWebhookContext(context.Background(), carrier, trackingNumber, metadata)
}

// RegisterTrackingWebhookContext is RegisterTrackingWebhook with cancellation and deadline support.
func (c *Client) RegisterTrackingWebhookContext(ctx context.Context, carrier, trackingNumber, metadata string) (*models.TrackingStatus, error) {
	return c.RegisterTrackingContext(ctx, &models.TrackingStatusInput{
		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		Metadata:       metadata,
	})
}

// RegisterTracking registers a tracking number with optional package details.
// Updates go to the account's configured tracking webhooks. This method sends
// the request once and does not retry it.
func (c *Client) RegisterTracking(input *models.TrackingStatusInput) (*models.TrackingStatus, error) {
	return c.RegisterTrackingContext(context.Background(), input)
}

// RegisterTrackingContext is RegisterTracking with cancellation and deadline support.
func (c *Client) RegisterTrackingContext(ctx context.Context, input *models.TrackingStatusInput) (*models.TrackingStatus, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}
	if input.Carrier == "" {
		return nil, errors.New("Empty carrier")
	}
	if input.TrackingNumber == "" {
		return nil, errors.New("Empty tracking number")
	}

	output := &models.TrackingStatus{}
	err := c.doContext(ctx, http.MethodPost, "/tracks/", input, output)
	return output, err
}
