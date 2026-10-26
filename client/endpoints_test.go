package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corvallis3d/go-shippo/models"
)

func endpointTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c := NewClientWithOptions("synthetic-token", "2018-02-08", Options{HTTPClient: server.Client()})
	c.baseURL = server.URL
	return c
}

func TestMultiParcelPurchaseAndFilteredTransactions(t *testing.T) {
	var purchases, pages int32
	c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Shippo-API-Version") != "2018-02-08" {
			t.Error("explicit API version missing")
		}
		switch r.Method + " " + r.URL.Path {
		case "POST /shipments/":
			var body struct {
				Parcels []models.ParcelInput `json:"parcels"`
				Extra   struct {
					QRCodeRequested bool `json:"qr_code_requested"`
				} `json:"extra"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if len(body.Parcels) != 2 || !body.Extra.QRCodeRequested {
				t.Errorf("shipment body = %+v", body)
			}
			fmt.Fprint(w, `{"object_id":"shipment-1","status":"SUCCESS","parcels":[{"object_id":"box-1"},{"object_id":"box-2"}]}`)
		case "GET /shipments/shipment-1/rates/USD":
			fmt.Fprint(w, `{"next":null,"results":[{"object_id":"rate-1","amount":"10.00","currency":"USD","servicelevel":{"token":"fedex_ground"}}]}`)
		case "POST /transactions/":
			atomic.AddInt32(&purchases, 1)
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			if body["rate"] != "rate-1" || body["order"] != "order-1" || body["async"] != false {
				t.Errorf("purchase body = %+v", body)
			}
			if _, ok := body["servicelevel_token"]; ok {
				t.Error("rate purchase should omit empty instant-purchase service level")
			}
			fmt.Fprint(w, `{"object_id":"tx-1","status":"SUCCESS","rate":"rate-1","parcel":"box-1","tracking_number":"track-1","label_url":"https://labels.example/1.pdf"}`)
		case "GET /transactions/":
			atomic.AddInt32(&pages, 1)
			if r.URL.Query().Get("rate") != "rate-1" || r.URL.Query().Get("results") != "1" {
				t.Errorf("filter lost on page: %s", r.URL.RawQuery)
			}
			if r.URL.Query().Get("page") == "2" {
				fmt.Fprint(w, `{"next":null,"results":[{"object_id":"tx-2","status":"SUCCESS","rate":{"object_id":"rate-1","amount":"10.00","currency":"USD"},"parcel":"box-2","tracking_number":"track-2","label_url":"https://labels.example/2.pdf"}]}`)
			} else {
				// A relative next link may omit initial filters; they must survive.
				fmt.Fprint(w, `{"next":"?page=2","results":[{"object_id":"tx-1","status":"SUCCESS","rate":"rate-1","parcel":"box-1","tracking_number":"track-1","label_url":"https://labels.example/1.pdf"}]}`)
			}
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	ctx := context.Background()
	input := &models.ShipmentInput{
		AddressFrom: "from-id", AddressTo: "to-id",
		Parcels: []*models.ParcelInput{
			{Length: "5", Width: "5", Height: "5", DistanceUnit: "in", Weight: "2", MassUnit: "lb"},
			{Length: "10", Width: "10", Height: "10", DistanceUnit: "in", Weight: "3", MassUnit: "lb"},
		},
		Extra: &models.ShipmentExtra{QRCodeRequested: true},
	}
	shipment, err := c.CreateShipmentContext(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if !input.ShipmentDate.IsZero() {
		t.Error("CreateShipment mutated caller input")
	}
	rates, err := c.GetShippingRatesContext(ctx, shipment.ObjectID, "USD")
	if err != nil || len(rates) != 1 {
		t.Fatalf("rates = %+v, error = %v", rates, err)
	}
	transaction, err := c.PurchaseShippingLabelContext(ctx, &models.TransactionInput{
		Rate: rates[0].ObjectID, Order: "order-1", LabelFileType: models.LabelFileTypePDF,
	})
	if err != nil || transaction.Parcel != "box-1" {
		t.Fatalf("purchase = %+v, error = %v", transaction, err)
	}
	labels, err := c.ListTransactionsContext(ctx, &models.TransactionListOptions{Rate: rates[0].ObjectID, Results: 1})
	if err != nil || len(labels) != 2 {
		t.Fatalf("labels = %+v, error = %v", labels, err)
	}
	for i, label := range labels {
		if label.Parcel != fmt.Sprintf("box-%d", i+1) || label.TrackingNumber != fmt.Sprintf("track-%d", i+1) || label.LabelURL == "" || label.Rate != "rate-1" {
			t.Errorf("parcel label %d = %+v", i, label)
		}
	}
	if atomic.LoadInt32(&purchases) != 1 || atomic.LoadInt32(&pages) != 2 {
		t.Fatalf("purchases = %d, pages = %d", purchases, pages)
	}
}

func TestLegacyTransactionListingAndEncodedFilter(t *testing.T) {
	var requests int32
	c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transactions/" || r.Method != http.MethodGet {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
		}
		if n := atomic.AddInt32(&requests, 1); n == 1 {
			if r.URL.Query().Get("rate") != "rate &?+1" {
				t.Errorf("unescaped filter: %q", r.URL.RawQuery)
			}
		} else if r.URL.Query().Get("rate") != "" {
			t.Error("legacy listing unexpectedly filtered")
		}
		fmt.Fprint(w, `{"next":null,"results":[{"object_id":"tx-1","rate":"rate-1"}]}`)
	})
	if _, err := c.ListTransactions(&models.TransactionListOptions{Rate: "rate &?+1"}); err != nil {
		t.Fatal(err)
	}
	if list, err := c.ListAllTransactions(); err != nil || len(list) != 1 {
		t.Fatalf("legacy listing = %+v, error = %v", list, err)
	}
	for _, results := range []int{-1, 101} {
		if _, err := c.ListTransactions(&models.TransactionListOptions{Results: results}); err == nil {
			t.Errorf("results %d accepted", results)
		}
	}
	if atomic.LoadInt32(&requests) != 2 {
		t.Error("invalid listing options sent a request")
	}
}

func TestOrderFiltersAndLegacyHelpers(t *testing.T) {
	start := time.Date(2026, 10, 26, 9, 0, 0, 0, time.FixedZone("offset", 2*60*60))
	c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "POST /orders/":
			var body models.OrderInput
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.OrderNumber != "order-1" {
				t.Errorf("order body = %+v, error = %v", body, err)
			}
			fmt.Fprint(w, `{"object_id":"order-id","order_number":"order-1","order_status":"PAID"}`)
		case "GET /orders/order-id":
			fmt.Fprint(w, `{"object_id":"order-id","order_status":"PAID"}`)
		case "GET /orders/":
			q := r.URL.Query()
			if q.Get("shop_app") != "Shippo & app" || q.Get("start_date") != "2026-10-26T07:00:00Z" || q.Get("end_date") != "2026-10-27T07:00:00Z" || q.Get("page") != "2" || q.Get("results") != "10" || !reflect.DeepEqual(q["order_status[]"], []string{"PAID", "PARTIALLY_FULFILLED"}) {
				t.Errorf("query = %v", q)
			}
			fmt.Fprint(w, `{"next":"?page=3","previous":"?page=1","results":[{"object_id":"order-id"}]}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	order, err := c.CreateOrder(&models.OrderInput{OrderNumber: "order-1", PlacedAt: start})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetOrder(order.ObjectID); err != nil {
		t.Fatal(err)
	}
	list, err := c.ListOrders(&models.OrderListOptions{
		ShopApp: "Shippo & app", StartDate: start, EndDate: start.Add(24 * time.Hour),
		Page: 2, Results: 10, OrderStatus: []string{"PAID", "PARTIALLY_FULFILLED"},
	})
	if err != nil || len(list.Results) != 1 || list.Next != "?page=3" {
		t.Fatalf("orders = %+v, error = %v", list, err)
	}
}

func TestAddressValidationAndRatesWithoutCurrency(t *testing.T) {
	c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "GET /addresses/address-id/validate":
			fmt.Fprint(w, `{"object_id":"address-id","validation_results":{"is_valid":true,"messages":[]}}`)
		case "GET /shipments/shipment-id/rates/":
			fmt.Fprint(w, `{"next":null,"results":[{"object_id":"rate-id"}]}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	address, err := c.ValidateAddress("address-id")
	if err != nil || address.ObjectID != "address-id" {
		t.Fatalf("address = %+v, error = %v", address, err)
	}
	if rates, err := c.GetShippingRates("shipment-id", ""); err != nil || len(rates) != 1 {
		t.Fatalf("rates = %+v, error = %v", rates, err)
	}
}
