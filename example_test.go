package shippo_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/corvallis3d/go-shippo"
	"github.com/corvallis3d/go-shippo/models"
)

func ExampleNewClientWithVersion() {
	token := os.Getenv("PRIVATE_TOKEN")
	if token == "" {
		log.Fatal("Set PRIVATE_TOKEN to your Shippo API token")
	}
	c := shippo.NewClientWithVersion(token, shippo.APIVersion20180208)
	shipment, err := c.CreateShipmentContext(context.Background(), &models.ShipmentInput{
		AddressFrom: &models.AddressInput{
			Name: "Mr. Hippo", Street1: "215 Clayton St.",
			City: "San Francisco", State: "CA", Zip: "94117", Country: "US",
		},
		AddressTo: &models.AddressInput{
			Name: "Mrs. Hippo", Street1: "965 Mission St.",
			City: "San Francisco", State: "CA", Zip: "94105", Country: "US",
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
	if len(shipment.Rates) == 0 {
		log.Fatalf("No rates returned; shipment status: %s", shipment.Status)
	}
	for _, rate := range shipment.Rates {
		fmt.Printf("%s: %s %s (%s)\n", rate.ObjectID, rate.Amount, rate.Currency, rate.Provider)
	}
}
