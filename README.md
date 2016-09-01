# Shippo API Golang Wrapper

Documentation: https://godoc.org/github.com/d5/go-shippo/client

## Example
_PLEASE NOTE: if you uncomment `purchaseShippingLabel()` line, Shippo's gonna **charge** you for real :)_
```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/d5/go-shippo/client"
	"github.com/d5/go-shippo/models"
)

func main() {
	privateToken := os.Getenv("PRIVATE_TOKEN")
	if privateToken == "" {
		panic(errors.New("Please set $PRIVATE_TOKEN with your Shippo API private token."))
	}

	// create a Shippo Client instance
	c := client.NewClient(privateToken)

	// create shipment
	shipment := createShipment(c)

	// purchase shipping label
	// NOTE: Uncomment the following line if you want to test the actual purchase.
	// purchaseShippingLabel(c, shipment)
}

func createShipment(c *client.Client) *models.ShipmentOutput {
	// create a sending address
	addressFromInput := &models.AddressInput{
		ObjectPurpose: models.ObjectPurposePurchase,
		Name:          "Mr. Hippo",
		Street1:       "215 Clayton St.",
		City:          "San Francisco",
		State:         "CA",
		Zip:           "94117",
		Country:       "US",
		Phone:         "+1 555 341 9393",
		Email:         "support@goshippo.com",
	}
	addressFrom, err := c.CreateAddress(addressFromInput)
	if err != nil {
		panic(err)
	}

	// create a receiving address
	addressToInput := &models.AddressInput{
		ObjectPurpose: models.ObjectPurposePurchase,
		Name:          "Mrs. Hippo",
		Street1:       "965 Mission St.",
		City:          "San Francisco",
		State:         "CA",
		Zip:           "94105",
		Country:       "US",
		Phone:         "+1 555 341 9393",
		Email:         "support@goshippo.com",
	}
	addressTo, err := c.CreateAddress(addressToInput)
	if err != nil {
		panic(err)
	}

	// create a parcel
	parcelInput := &models.ParcelInput{
		Length:       "5",
		Width:        "5",
		Height:       "5",
		DistanceUnit: models.DistanceUnitInch,
		Weight:       "2",
		MassUnit:     models.MassUnitPound,
	}
	parcel, err := c.CreateParcel(parcelInput)
	if err != nil {
		panic(err)
	}

	// create a shipment
	shipmentInput := &models.ShipmentInput{
		ObjectPurpose: models.ObjectPurposePurchase,
		AddressFrom:   addressFrom.ObjectID,
		AddressTo:     addressTo.ObjectID,
		Parcel:        parcel.ObjectID,
		Async:         false,
	}
	shipment, err := c.CreateShipment(shipmentInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shipment:\n%s\n", dump(shipment))

	return shipment
}

func purchaseShippingLabel(c *client.Client, shipment *models.ShipmentOutput) {
	transactionInput := &models.TransactionInput{
		Rate:          shipment.RatesList[len(shipment.RatesList)-1].ObjectID,
		LabelFileType: models.LabelFileTypePDF,
		Async:         false,
	}
	transaction, err := c.PurchaseShippingLabel(transactionInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction:\n%s\n", dump(transaction))
}

func dump(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}
```
To run:
```bash
PRIVATE_TOKEN=<your_api_token> go run *.go
```
