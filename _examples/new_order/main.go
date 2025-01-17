package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/corvallis3d/go-shippo"
	"github.com/corvallis3d/go-shippo/client"
	"github.com/corvallis3d/go-shippo/models"
)

func main() {
	privateToken := os.Getenv("PRIVATE_TOKEN")
	if privateToken == "" {
		panic(errors.New("Please set $PRIVATE_TOKEN with your Shippo API private token."))
	}

	// create a Shippo Client instance
	c := shippo.NewClient(privateToken)

	// create order
	createOrder(c)
}

func createOrder(c *client.Client) *models.OrderResponse {
	// create a sending address
	toAddressInput := &models.AddressInput{
		Name:     "Shawn Ippotle",
		Company:  "Shippo",
		Street1:  "215 Clayton St.",
		City:     "San Francisco",
		State:    "CA",
		Zip:      "94117",
		Country:  "US",
		Email:    "shippotle@goshippo.com",
		Validate: true,
	}
	toAddress, err := c.CreateAddress(toAddressInput)
	if err != nil {
		panic(err)
	}

	item1 := &models.LineItem{
		ObjectID:   "ID#123456",
		Title:      "Print1",
		Sku:        "32456312",
		Quantity:   2,
		TotalPrice: "23.70",
		Weight:     "5",
		Currency:   "USD",
		WeightUnit: "lb",
	}

	// create a shipment
	orderInput := &models.OrderInput{
		ToAddress:            toAddress.ObjectID,
		LineItems:            []models.LineItem{*item1},
		OrderNumber:          "231442",
		Currency:             "USD",
		ShippingCostCurrency: "USD",
		Weight:               "4",
		WeightUnit:           "lb",
	}
	order, err := c.CreateOrder(orderInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Order:\n%s\n", dump(order))

	return order
}

func dump(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}
