package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/corvallis3d/go-shippo/models"
)

// PurchaseShippingLabel creates a new transaction object and purchases the shipping label for the provided rate.
func (c *Client) PurchaseShippingLabel(input *models.TransactionInput) (*models.Transaction, error) {
	return c.PurchaseShippingLabelContext(context.Background(), input)
}

// PurchaseShippingLabelContext is PurchaseShippingLabel with cancellation and deadline support.
func (c *Client) PurchaseShippingLabelContext(ctx context.Context, input *models.TransactionInput) (*models.Transaction, error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	output := &models.Transaction{}
	err := c.doContext(ctx, http.MethodPost, "/transactions/", input, output)
	return output, err
}

// RetrieveTransaction retrieves an existing transaction by object id.
func (c *Client) RetrieveTransaction(objectID string) (*models.Transaction, error) {
	return c.RetrieveTransactionContext(context.Background(), objectID)
}

// RetrieveTransactionContext is RetrieveTransaction with cancellation and deadline support.
func (c *Client) RetrieveTransactionContext(ctx context.Context, objectID string) (*models.Transaction, error) {
	if objectID == "" {
		return nil, errors.New("Empty object ID")
	}

	output := &models.Transaction{}
	err := c.doContext(ctx, http.MethodGet, "/transactions/"+objectID, nil, output)
	return output, err
}

// ListAllTransactions lists all transaction objects.
func (c *Client) ListAllTransactions() ([]*models.Transaction, error) {
	return c.ListAllTransactionsContext(context.Background())
}

// ListAllTransactionsContext is ListAllTransactions with cancellation and deadline support.
func (c *Client) ListAllTransactionsContext(ctx context.Context) ([]*models.Transaction, error) {
	return c.ListTransactionsContext(ctx, nil)
}

// ListTransactions lists all transactions matching options, following every page.
// A rate filter retrieves the individual parcel labels of a multi-piece shipment.
func (c *Client) ListTransactions(options *models.TransactionListOptions) ([]*models.Transaction, error) {
	return c.ListTransactionsContext(context.Background(), options)
}

// ListTransactionsContext is ListTransactions with cancellation and deadline support.
// If a later page fails, it returns the transactions already read along with the error.
func (c *Client) ListTransactionsContext(ctx context.Context, options *models.TransactionListOptions) ([]*models.Transaction, error) {
	path := "/transactions/"
	query := url.Values{}
	if options != nil {
		if options.Rate != "" {
			query.Set("rate", options.Rate)
		}
		if options.Results < 0 || options.Results > 100 {
			return nil, errors.New("results must be between 1 and 100, or zero for the default")
		}
		if options.Results > 0 {
			query.Set("results", strconv.Itoa(options.Results))
		}
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	list := []*models.Transaction{}
	err := c.doListContext(ctx, http.MethodGet, path, nil, func(v json.RawMessage) error {
		item := &models.Transaction{}
		if err := json.Unmarshal(v, item); err != nil {
			return err
		}

		list = append(list, item)
		return nil
	})
	return list, err
}
