package client

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	apierrors "github.com/corvallis3d/go-shippo/errors"
	"github.com/corvallis3d/go-shippo/models"
)

// October 8, 2026 added per-transaction rejection reasons to manifest errors.
func TestManifestRejectionDetailsPreserved(t *testing.T) {
	const response = `{"detail":"Invalid transactions transaction-1,transaction-2 for manifest request","transactions":{"transaction-1":{"reason":"already_manifested","manifest_object_id":"manifest-1","manifest_status":"SUCCESS"},"transaction-2":{"reason":"outside_submission_date_window","shipment_submission_date":"2026-10-07T12:00:00Z","window_start":"2026-10-08T00:00:00Z","window_end":"2026-10-08T23:59:59Z"}}}`
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, response)
	}, Options{})
	_, err := c.CreateManifest(&models.ManifestInput{
		CarrierAccount: "carrier-1", AddressFrom: "address-1",
		ShipmentDate: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
		Transactions: []string{"transaction-1", "transaction-2"},
	})
	var apiErr *apierrors.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound || string(apiErr.ResponseBody) != response {
		t.Fatalf("manifest rejection details lost: %v", err)
	}
}
