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

	// get order
	getOrder(c)
}

func getOrder(c *client.Client) *models.OrderResponse {
	// order by object_id
	orderId := "1234567"

	order, err := c.GetOrder(orderId)
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
