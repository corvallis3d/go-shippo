# Shippo API Golang Wrapper

Documentation: https://godoc.org/github.com/corvallis3d/go-shippo/client

## Client configuration and API coverage

This is a thin wrapper around Shippo's models and endpoints. Share one client
across requests; every endpoint has a `Context` variant. Existing methods and
constructors remain available.

```go
c := shippo.NewClientWithOptions(token, shippo.APIVersion20180208, client.Options{
    HTTPClient: &http.Client{},
    Timeout:    45 * time.Second,
})
ctx, cancel := context.WithTimeout(parentContext, 2*time.Minute)
defer cancel()

rates, err := c.GetShippingRatesContext(ctx, shipmentID, "USD")
// Handle err and let the application select a rate from rates.
purchase, err := c.PurchaseShippingLabelContext(ctx, &models.TransactionInput{
    Rate: selectedRateID, Order: orderID,
    LabelFileType: models.LabelFileTypePDF, Async: false,
})
// Handle err and purchase.Status before continuing; never retry an uncertain purchase.
labels, err := c.ListTransactionsContext(ctx, &models.TransactionListOptions{
    Rate: selectedRateID,
})
// Check err before treating labels as complete. Each transaction provides
// Parcel, Status, LabelURL, TrackingNumber, and (when available) QRCodeURL.
```

See [the complete, compile-checked example](example_test.go) for selection and
error handling. Create multi-piece shipments using the existing
`ShipmentInput.Parcels` list of parcel IDs or `[]*models.ParcelInput`.
Request QR codes using `ShipmentInput.Extra.QRCodeRequested`; QR availability
depends on the carrier. The response URL is `Transaction.QRCodeURL`.

### HTTP behavior and compatibility

- `NewClient(token)` still omits `Shippo-API-Version`, using the account default.
  `NewClientWithVersion` keeps its behavior. Use `APIVersion20180208` explicitly
  for the current core reference; no account settings are changed by this library.
- Each HTTP request, including its response body, has a **30-second default timeout**.
  A positive `Options.Timeout` overrides the supplied HTTP client's timeout.
  Otherwise a positive supplied timeout is retained; zero/negative values use
  the default. A context can end a request sooner and bound a whole paginated list.
- The supplied `http.Client` is copied, retaining its transport/pool. Configure
  transports before sharing the client. Injected transports must honor contexts
  and must not retry writes. The wrapper performs no automatic retries.
- All redirects return an `*errors.APIError` with the 3xx status, including
  same-origin redirects. Pagination accepts only the API origin (same scheme,
  host and port, with no userinfo), supporting absolute and relative next links.
  Connections are reused. Trace logs contain only methods and status codes.
- HTTP errors retain `APIError.Status` and `ResponseBody`. Use `errors.Is` for
  cancellation/deadlines and `errors.As` for HTTP/network errors. A timeout,
  connection loss, or response decoding failure after purchase may have an
  uncertain remote outcome. The application owns reconciliation and checkpoints.
- `ListTransactions` follows every page and optionally filters by rate.
  Existing `ListAll*` helpers use the same paginator. Errors can accompany partial
  list results. `ListOrders` retains its existing **one-page** behavior.
- Existing Go names remain, including `ServerLevelToken`,
  `APIVesrion20161025` (with a correctly spelled alias), and `Rate.Days`.
  Current responses populate `Rate.EstimatedDays`. A transaction's `Rate` stays
  an ID when the API returns an expanded object; details are in `RateDetails`.
- Wire corrections include `servicelevel_token`, `PDF_4x6`, omission of an unset
  optional label format/date or template parcel dimensions, and removal of
  read-only fields from order addresses. Existing exported field types remain.
  Use keyed struct literals when constructing models that gain optional fields.

### Official API audit (2026-10-26)

Verified the [documentation index](https://docs.goshippo.com/llms.txt),
[current reference](https://docs.goshippo.com/api-reference/overview),
[official OpenAPI schema](https://docs.goshippo.com/spec/shippoapi/public-api.yaml),
[release notes through October 20, 2026](https://docs.goshippo.com/api-concepts/api-change-log#oct-20-2026),
[versioning](https://docs.goshippo.com/api-concepts/api-versioning), and
[multi-piece guide](https://docs.goshippo.com/shipments/multi-piece-shipment).
The core reference remains **2018-02-08**, served at
`https://api.goshippo.com` (the old `/v1` prefix was removed).
The version date identifies the API compatibility contract, not the age of its
features. Shippo continues adding endpoints and optional fields under that same
version. This wrapper is checked against the current schema, guides, and changelog.
The public [documentation MCP](https://docs.goshippo.com/mcp) was also queried for
schema fields and the latest release notes. It exposes documentation search and
the published OpenAPI specs without shipping-account access. Its core schema still
omits some guide-documented additions, such as tracking `include_package_details`,
so the changelog and linked guides are checked alongside it. Shippo's separate
account MCP and AI knowledge pack provide operational tools/workflow guidance;
their tool catalog is not a complete REST API coverage checklist.

| Area | Coverage and corrections |
| --- | --- |
| Addresses | Create, retrieve, list; added existing-address validation and fixed optional state tag. |
| Orders | Create, retrieve, filtered one-page listing; safe query encoding, response count and clean address request shapes. |
| Parcels | Create, retrieve, list; template parcels can omit dimensions. Existing multi-parcel forms retained. |
| Shipments/rates | Create/retrieve/list shipments and retrieve/list rates; optional rate currency, current estimates, shipment extras and metadata; customs ID/object response compatibility. |
| Transactions/labels | Purchase/retrieve/list; order association, parcel IDs, rate-filtered pagination, expanded rate decoding, service-level and label-format fixes. |
| QR codes | Request on shipment extras; URL on transaction, per the [QR guide](https://docs.goshippo.com/shipments/qr-code). |
| Tracking | Retrieve/register; pre-transit status, original ETA, transaction association, event substatus, and optional package weight/dimensions. |
| Customs | Create/retrieve/list items and declarations; importer reference, product identifiers, exporter tax/EORI, duties payor, B13A and VAT fields. |
| Refunds | Existing create/retrieve/list shapes verified with synthetic fixtures. |
| Batches/manifests | Existing endpoints retained; nested batch shipment requests, batch shipment status, manifest status/errors and raw rejection details preserved. |
| Carrier accounts | Existing endpoints retained; carrier metadata/service levels and authentication/verification state exposed. |

Recent changelog coverage in the existing wrapper:

| Release entry | Handling |
| --- | --- |
| Oct 20, 2026 | Added `healthcare_handling`, `reference_3`, `reference_4`; new service-level tokens remain unrestricted strings. |
| Oct 1/14, 2026 | Customs product identifiers and recipient-business fields; missing-field warnings remain in shipment/transaction messages. Carrier enforcement is left to Shippo. |
| Oct 8, 2026 | Manifest 404 rejection details, including per-transaction reason objects, remain intact in `APIError.ResponseBody`. |
| Sep 23, 2026 | Tracking event object IDs are retained. |
| Jun 24, 2026 | FedEx service-level `terms`, including One Rate, are retained alongside the token. |
| Apr 6, 2026 | Shippo's 390-day shipment/rate retention is a server policy; HTTP errors are surfaced without retries. |
| Mar 5, 2026 | `RegisterTracking[Context]` accepts `TrackingStatusInput.IncludePackageDetails`; response `Weight` and `Dimensions` are available. Legacy registration methods still work. |
| Mar 4, 2026 | Carrier account `ObjectInfo.Authentication` and `ObjectInfo.Verification.Policy` are available. |

Earlier compatible model additions include FedEx/UPS label references, dangerous
goods, optional insurance provider for XCover, additional label formats, customs
tax identifiers, and order validation/error responses. No carrier/service enum
validation blocks newly introduced tokens.

The [Orders API](https://docs.goshippo.com/orders/orders) documents create/get/list,
not a general update/delete operation. Shippo controls order status after creation.
Order synchronization policy remains in the application.

Deferred API areas: order packing slips; batch creation; carrier OAuth/registration;
pickups; parcel-template discovery; service groups; checkout rates/settings; managed
accounts; webhook management; embedded JWT; Addresses v2; reporting; estimates;
Shippo Automations beta; and exhaustive carrier-specific extras/advanced customs.
Existing ancillary endpoints remain supported. No business workflows, UI models,
carrier-specific parcel limits, polling loops, or purchase recovery were added.

Validation uses synthetic JSON and local HTTP mocks only: request/response shapes,
multi-page parcel labels, legacy calls, cancellation/timeouts, HTTP errors, connection
reuse, redirect/pagination safety, and single-attempt purchases. Run:

```sh
gofmt -l .
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Suggested release: **v1.7.0** (additive client capabilities).
Suggested commit: `Add context-aware Shippo client and multi-package label support`.

## Basic example using legacy methods

```go
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

	// create shipment
	shipment := createShipment(c)

	// purchase shipping label
	purchaseShippingLabel(c, shipment)
}

func createShipment(c *client.Client) *models.Shipment {
	// create a sending address
	addressFromInput := &models.AddressInput{
		Name:    "Mr. Hippo",
		Street1: "215 Clayton St.",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94117",
		Country: "US",
		Phone:   "+1 555 341 9393",
		Email:   "support@goshippo.com",
	}
	addressFrom, err := c.CreateAddress(addressFromInput)
	if err != nil {
		panic(err)
	}

	// create a receiving address
	addressToInput := &models.AddressInput{
		Name:    "Mrs. Hippo",
		Street1: "965 Mission St.",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94105",
		Country: "US",
		Phone:   "+1 555 341 9393",
		Email:   "support@goshippo.com",
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
		AddressFrom: addressFrom.ObjectID,
		AddressTo:   addressTo.ObjectID,
		Parcels:     []string{parcel.ObjectID},
		Async:       false,
	}
	shipment, err := c.CreateShipment(shipmentInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shipment:\n%s\n", dump(shipment))

	return shipment
}

func purchaseShippingLabel(c *client.Client, shipment *models.Shipment) {
	transactionInput := &models.TransactionInput{
		Rate:          shipment.Rates[len(shipment.Rates)-1].ObjectID,
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

See [more examples](https://github.com/corvallis3d/go-shippo/tree/master/_examples).
