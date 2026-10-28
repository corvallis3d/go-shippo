package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/corvallis3d/go-shippo/models"
)

// RetrieveBatch retrieves a batch, with up to 100 shipments per page.
// A nonzero page selects a page; objectResultsFilter selects shipment results.
func (c *Client) RetrieveBatch(objectID string, page uint, objectResultsFilter string) (*models.Batch, error) {
	return c.RetrieveBatchContext(context.Background(), objectID, page, objectResultsFilter)
}

// RetrieveBatchContext is RetrieveBatch with cancellation and deadline support.
func (c *Client) RetrieveBatchContext(ctx context.Context, objectID string, page uint, objectResultsFilter string) (*models.Batch, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	url := "/batches/" + objectID
	qs := []string{}
	if page > 0 {
		qs = append(qs, fmt.Sprintf("page=%d", page))
	}
	if objectResultsFilter != "" {
		qs = append(qs, "object_results="+objectResultsFilter)
	}
	if len(qs) > 0 {
		url += "?" + strings.Join(qs, "&")
	}

	output := &models.Batch{}
	err := c.doContext(ctx, http.MethodGet, url, nil, output)
	return output, err
}

// AddBatchShipmentsToBatch adds batch shipment(s) to an existing Batch.
func (c *Client) AddBatchShipmentsToBatch(objectID string, batchShipments []*models.BatchShipmentInput) (*models.Batch, error) {
	return c.AddBatchShipmentsToBatchContext(context.Background(), objectID, batchShipments)
}

// AddBatchShipmentsToBatchContext is AddBatchShipmentsToBatch with cancellation and deadline support.
func (c *Client) AddBatchShipmentsToBatchContext(ctx context.Context, objectID string, batchShipments []*models.BatchShipmentInput) (*models.Batch, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}
	if batchShipments == nil {
		return nil, errors.New("nil batch shipments")
	}

	output := &models.Batch{}
	err := c.doContext(ctx, http.MethodPost, "/batches/"+objectID+"/add_shipments", &batchShipments, output)
	return output, err
}

// RemoveBatchShipmentsFromBatch removes batch shipment(s) from an existing Batch.
func (c *Client) RemoveBatchShipmentsFromBatch(objectID string, batchShipmentIDs []string) (*models.Batch, error) {
	return c.RemoveBatchShipmentsFromBatchContext(context.Background(), objectID, batchShipmentIDs)
}

// RemoveBatchShipmentsFromBatchContext is RemoveBatchShipmentsFromBatch with cancellation and deadline support.
func (c *Client) RemoveBatchShipmentsFromBatchContext(ctx context.Context, objectID string, batchShipmentIDs []string) (*models.Batch, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}
	if batchShipmentIDs == nil {
		return nil, errors.New("nil batch shipment IDs")
	}

	output := &models.Batch{}
	err := c.doContext(ctx, http.MethodPost, "/batches/"+objectID+"/remove_shipments", &batchShipmentIDs, output)
	return output, err
}

// PurchaseBatch purchases a batch whose Status is VALID.
// The batch moves through PURCHASING to PURCHASED asynchronously.
func (c *Client) PurchaseBatch(objectID string) (*models.Batch, error) {
	return c.PurchaseBatchContext(context.Background(), objectID)
}

// PurchaseBatchContext is PurchaseBatch with cancellation and deadline support.
func (c *Client) PurchaseBatchContext(ctx context.Context, objectID string) (*models.Batch, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Batch{}
	err := c.doContext(ctx, http.MethodPost, "/batches/"+objectID+"/purchase", nil, output)
	return output, err
}
