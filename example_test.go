package shippo_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	shippo "github.com/corvallis3d/go-shippo"
	"github.com/corvallis3d/go-shippo/client"
	"github.com/corvallis3d/go-shippo/models"
)

// This example is compiled but never run by go test (there is no Output directive).
func ExampleNewClientWithOptions() {
	// Construct once at application startup and share across operations.
	c := shippo.NewClientWithOptions(os.Getenv("SHIPPO_TOKEN"), shippo.APIVersion20180208, client.Options{
		HTTPClient: &http.Client{},
		Timeout:    45 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// The shipment may contain multiple parcels. Its creation response also has rates.
	rates, err := c.GetShippingRatesContext(ctx, "shipment-object-id", "USD")
	if err != nil {
		return // Report the error to the application.
	}
	// The application selects a returned rate; this ID stands for that selection.
	selectedRateID := "selected-rate-object-id"
	var selected *models.Rate
	for _, rate := range rates {
		if rate.ObjectID == selectedRateID {
			selected = rate
			break
		}
	}
	if selected == nil {
		return
	}

	purchase, err := c.PurchaseShippingLabelContext(ctx, &models.TransactionInput{
		Rate: selected.ObjectID, Order: "order-object-id",
		LabelFileType: models.LabelFileTypePDF, Async: false,
	})
	if err != nil {
		return // The purchase may have succeeded remotely. Do not repurchase here.
	}
	if purchase.Status != models.StatusSuccess {
		return // Give the status/messages and transaction ID to the application's workflow.
	}

	// Follow every page for the purchased rate, including each parcel's transaction.
	labels, err := c.ListTransactionsContext(ctx, &models.TransactionListOptions{Rate: selected.ObjectID})
	if err != nil {
		return // Results may be partial; resume/reconcile reads in the application.
	}
	for _, label := range labels {
		fmt.Println(label.Parcel, label.Status, label.LabelURL, label.TrackingNumber)
	}
}
