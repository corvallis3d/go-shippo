# go-shippo

A small, dependency-free Go client for [Shippo](https://docs.goshippo.com/).
Maintained fork of [coldbrewcloud/go-shippo](https://github.com/coldbrewcloud/go-shippo).

[Go documentation](https://pkg.go.dev/github.com/corvallis3d/go-shippo/client) ·
[Shippo API reference](https://docs.goshippo.com/api-reference/overview)

## Install

```sh
go get github.com/corvallis3d/go-shippo
```

## Usage

Set `PRIVATE_TOKEN` to a Shippo test token to try the examples. Create a shipment
and inspect its rates:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/corvallis3d/go-shippo"
    "github.com/corvallis3d/go-shippo/models"
)

func main() {
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
```

See [first shipment](_examples/first_shipment/main.go) for label purchase using
`SERVICELEVEL_TOKEN`, or [carrier accounts](_examples/carrier_accounts/main.go)
for an already-connected `CARRIER_ACCOUNT_ID`. Other examples cover
[address validation](_examples/address_validation/main.go) and
[orders](_examples/new_order/main.go).

## Client behavior

- `NewClient` uses the account's default API version. The example explicitly
  selects Shippo's `2018-02-08` core API.
- Reuse a client across requests. Every endpoint has a `Context` variant. Requests
  default to a 30-second timeout; `NewClientWithOptions` accepts a custom
  `http.Client` and timeout. Redirects are disabled.
- HTTP errors expose `errors.APIError.Status` and `ResponseBody`. A label purchase
  also needs `Status == models.StatusSuccess`. Purchases are not automatically
  retried; a timeout may leave the outcome unknown.
- `ListAll*` and `ListTransactions` follow all pages and may return partial results
  with an error. `ListOrders` returns one page.
- For multi-package labels, use `ListTransactions` with
  `models.TransactionListOptions{Rate: rateID}`. `Transaction.Rate` holds the ID;
  `RateDetails` holds expanded rate data when returned by Shippo.

## Development

```sh
go test -race ./...
go vet ./...
go build ./...
```

Tests use local HTTP servers and JSON fixtures; no Shippo credentials are required.

## License

[MIT](LICENSE), retaining the original copyright notice.
