package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/corvallis3d/go-shippo"
	"github.com/corvallis3d/go-shippo/models"
)

func main() {
	privateToken := os.Getenv("PRIVATE_TOKEN")
	serviceToken := os.Getenv("SERVICELEVEL_TOKEN") // For example, usps_priority.
	if privateToken == "" || serviceToken == "" {
		log.Fatal("Set PRIVATE_TOKEN and SERVICELEVEL_TOKEN; use a Shippo test token for this example.")
	}

	c := shippo.NewClientWithVersion(privateToken, shippo.APIVersion20180208)
	shipment, err := c.CreateShipment(&models.ShipmentInput{
		AddressFrom: &models.AddressInput{
			Name: "Mr. Hippo", Street1: "215 Clayton St.",
			City: "San Francisco", State: "CA", Zip: "94117", Country: "US",
			Phone: "+1 555 341 9393", Email: "support@goshippo.com",
		},
		AddressTo: &models.AddressInput{
			Name: "Mrs. Hippo", Street1: "965 Mission St.",
			City: "San Francisco", State: "CA", Zip: "94105", Country: "US",
			Phone: "+1 555 341 9393", Email: "support@goshippo.com",
		},
		Parcels: []*models.ParcelInput{{
			Length: "5", Width: "5", Height: "5", DistanceUnit: models.DistanceUnitInch,
			Weight: "2", MassUnit: models.MassUnitPound,
		}},
		Async: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	var selected *models.Rate
	for _, rate := range shipment.Rates {
		if rate == nil || rate.ServiceLevel == nil || rate.ServiceLevel.Token != serviceToken {
			continue
		}
		if selected != nil {
			log.Fatal("Multiple rates match this service; choose a specific carrier account before purchasing.")
		}
		selected = rate
	}
	if selected == nil {
		messages, _ := json.Marshal(shipment.Messages)
		log.Fatalf("No rate for %q; shipment status=%s messages=%s", serviceToken, shipment.Status, messages)
	}

	transaction, err := c.PurchaseShippingLabel(&models.TransactionInput{
		Rate: selected.ObjectID, LabelFileType: models.LabelFileTypePDF, Async: false,
	})
	if err != nil {
		log.Fatalf("Purchase outcome may be unknown; reconcile before retrying: %v", err)
	}
	if transaction.Status != models.StatusSuccess {
		messages, _ := json.Marshal(transaction.Messages)
		log.Fatalf("Transaction %s: status=%s messages=%s", transaction.ObjectID, transaction.Status, messages)
	}
	fmt.Println(transaction.LabelURL)
}
