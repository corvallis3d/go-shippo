package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	shippoerrors "github.com/corvallis3d/go-shippo/errors"
	"github.com/corvallis3d/go-shippo/models"
)

func TestRegisterTrackingPackageDetailsAndLegacyHelpers(t *testing.T) {
	for _, mode := range []string{"context", "background", "legacy", "legacy-context"} {
		t.Run(mode, func(t *testing.T) {
			var requests int32
			includeDetails := mode == "context" || mode == "background"
			c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				if r.Method != http.MethodPost || r.URL.Path != "/tracks/" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL)
				}
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body["carrier"] != "ups" || body["tracking_number"] != "track-1" || body["metadata"] != "order-1" {
					t.Errorf("tracking body = %+v", body)
				}
				if details, present := body["include_package_details"]; includeDetails {
					if !present || details != true {
						t.Error("package details option missing")
					}
				} else if present {
					t.Error("legacy registration should omit package details option")
				}
				fmt.Fprint(w, `{"carrier":"ups","tracking_number":"track-1","tracking_status":{"status":"TRANSIT"},"weight":[{"value":"2.2","unit":"LB"}],"dimensions":[{"length":"5","unit":"in"}]}`)
			})
			input := &models.TrackingStatusInput{
				Carrier: "ups", TrackingNumber: "track-1", Metadata: "order-1", IncludePackageDetails: true,
			}
			var status *models.TrackingStatus
			var err error
			switch mode {
			case "context":
				status, err = c.RegisterTrackingContext(context.Background(), input)
			case "background":
				status, err = c.RegisterTracking(input)
			case "legacy":
				status, err = c.RegisterTrackingWebhook("ups", "track-1", "order-1")
			case "legacy-context":
				status, err = c.RegisterTrackingWebhookContext(context.Background(), "ups", "track-1", "order-1")
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(status.Weight) != 1 || status.Weight[0].Value != "2.2" || len(status.Dimensions) != 1 || status.Dimensions[0].Length != "5" || status.Dimensions[0].Width != "" {
				t.Errorf("package details = %+v, %+v", status.Weight, status.Dimensions)
			}
			if atomic.LoadInt32(&requests) != 1 {
				t.Errorf("registration sent %d times", requests)
			}
		})
	}
}

func TestRegisterTrackingFailureIsNotRetried(t *testing.T) {
	var requests int32
	c := endpointTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `{"detail":"synthetic upstream error"}`)
	})
	for _, input := range []*models.TrackingStatusInput{nil, {}, {Carrier: "ups"}} {
		if _, err := c.RegisterTracking(input); err == nil {
			t.Errorf("invalid registration accepted: %+v", input)
		}
	}
	if atomic.LoadInt32(&requests) != 0 {
		t.Fatal("invalid registration sent a request")
	}
	_, err := c.RegisterTrackingContext(context.Background(), &models.TrackingStatusInput{Carrier: "ups", TrackingNumber: "track-1"})
	var apiErr *shippoerrors.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusBadGateway {
		t.Fatalf("registration error = %v", err)
	}
	if atomic.LoadInt32(&requests) != 1 {
		t.Errorf("failed registration sent %d times", requests)
	}
}
