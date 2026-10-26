package client

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	apierrors "github.com/corvallis3d/go-shippo/errors"
	"github.com/corvallis3d/go-shippo/models"
)

func localClient(t *testing.T, handler http.HandlerFunc, options Options) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	c := NewClientWithOptions("secret-test-token", "2018-02-08", options)
	c.baseURL = server.URL
	return c
}

func TestHTTPConfigurationAndConnectionReuse(t *testing.T) {
	var addresses []string
	var mu sync.Mutex
	supplied := &http.Client{Timeout: 2 * time.Second}
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		addresses = append(addresses, r.RemoteAddr)
		mu.Unlock()
		if r.URL.Path != "/addresses/address-id" ||
			r.Header.Get("Authorization") != "ShippoToken secret-test-token" ||
			r.Header.Get("Shippo-API-Version") != "2018-02-08" ||
			r.Close || r.Header.Get("Connection") != "" {
			t.Errorf("unexpected request configuration: path=%s close=%t", r.URL.Path, r.Close)
		}
		io.WriteString(w, "{}")
	}, Options{HTTPClient: supplied, Timeout: time.Second})
	if c.httpClient == supplied || supplied.Timeout != 2*time.Second || supplied.CheckRedirect != nil {
		t.Fatal("constructor mutated the supplied HTTP client")
	}
	if c.httpClient.Timeout != time.Second {
		t.Fatal("timeout override not used")
	}
	for i := 0; i < 2; i++ {
		if _, err := c.RetrieveAddress("address-id"); err != nil {
			t.Fatal(err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if len(addresses) != 2 || addresses[0] != addresses[1] {
		t.Fatalf("connection was not reused: %v", addresses)
	}
}

func TestTimeoutDefaultsAndVersion(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options Options
		want    time.Duration
	}{
		{"default", Options{}, DefaultTimeout},
		{"supplied", Options{HTTPClient: &http.Client{Timeout: time.Minute}}, time.Minute},
		{"negative", Options{HTTPClient: &http.Client{Timeout: -1}, Timeout: -1}, DefaultTimeout},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClientWithOptions("token", "", tt.options)
			if c.httpClient.Timeout != tt.want {
				t.Fatalf("timeout=%v want %v", c.httpClient.Timeout, tt.want)
			}
			r, err := c.createRequest(context.Background(), "GET", shippoAPIBaseURL+"/orders/", nil)
			if err != nil || r.Header.Get("Shippo-API-Version") != "" || r.URL.Path != "/orders/" {
				t.Fatalf("account-default version or current base URL changed: %v", err)
			}
		})
	}
}

func TestPurchaseCancellationAndTimeoutSendOnce(t *testing.T) {
	for _, mode := range []string{"cancel", "context deadline", "client timeout"} {
		t.Run(mode, func(t *testing.T) {
			var requests int32
			started := make(chan struct{})
			release := make(chan struct{})
			options := Options{Timeout: time.Second}
			if mode == "client timeout" {
				options.Timeout = 100 * time.Millisecond
			}
			c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				io.Copy(io.Discard, r.Body)
				close(started)
				<-release
			}, options)
			defer close(release)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if mode == "context deadline" {
				ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
			}
			done := make(chan error, 1)
			go func() {
				_, err := c.PurchaseShippingLabelContext(ctx, &models.TransactionInput{Rate: "rate-id"})
				done <- err
			}()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("purchase request did not arrive")
			}
			if mode == "cancel" {
				cancel()
			}
			err := <-done
			want := context.DeadlineExceeded
			if mode == "cancel" {
				want = context.Canceled
			}
			if !stderrors.Is(err, want) {
				t.Fatalf("lost cancellation/timeout cause: %v", err)
			}
			if mode != "cancel" {
				var timeout net.Error
				if !stderrors.As(err, &timeout) || !timeout.Timeout() {
					t.Fatalf("timeout is not discoverable: %v", err)
				}
			}
			if atomic.LoadInt32(&requests) != 1 {
				t.Fatalf("purchase attempted %d times", requests)
			}
		})
	}
}

func TestPreCanceledPurchaseIsNotSent(t *testing.T) {
	var requests int32
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
	}, Options{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.PurchaseShippingLabelContext(ctx, &models.TransactionInput{Rate: "rate"})
	if !stderrors.Is(err, context.Canceled) || atomic.LoadInt32(&requests) != 0 {
		t.Fatalf("pre-canceled purchase was sent or lost cause: %v", err)
	}
}

func TestHTTPFailuresAndBrokenPurchaseAreNotRetried(t *testing.T) {
	for _, status := range []int{400, 401, 429, 500, 503, 0} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			var requests int32
			c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				if status == 0 {
					conn, _, err := w.(http.Hijacker).Hijack()
					if err != nil {
						t.Error(err)
						return
					}
					conn.Close()
					return
				}
				w.WriteHeader(status)
				io.WriteString(w, "{\"detail\":\"synthetic failure\"}")
			}, Options{})
			_, err := c.PurchaseShippingLabel(&models.TransactionInput{Rate: "rate"})
			if err == nil || atomic.LoadInt32(&requests) != 1 {
				t.Fatalf("unexpected purchase outcome/count: %v / %d", err, requests)
			}
			if status != 0 {
				var apiErr *apierrors.APIError
				if !stderrors.As(err, &apiErr) || apiErr.Status != status ||
					!strings.Contains(string(apiErr.ResponseBody), "synthetic failure") {
					t.Fatalf("HTTP status/body not preserved: %v", err)
				}
			}
		})
	}
}

func TestRedirectsNeverForwardOrReplay(t *testing.T) {
	var targetRequests int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&targetRequests, 1)
	}))
	defer target.Close()
	for _, status := range []int{301, 302, 303, 307, 308} {
		for _, sameOrigin := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/same=%t", status, sameOrigin), func(t *testing.T) {
				var requests int32
				c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
					atomic.AddInt32(&requests, 1)
					location := target.URL + "/capture"
					if sameOrigin {
						location = "/capture"
					}
					w.Header().Set("Location", location)
					w.WriteHeader(status)
				}, Options{HTTPClient: &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
					t.Error("supplied redirect callback must not override safety policy")
					return nil
				}}})
				_, err := c.PurchaseShippingLabel(&models.TransactionInput{Rate: "rate"})
				var apiErr *apierrors.APIError
				if !stderrors.As(err, &apiErr) || apiErr.Status != status ||
					atomic.LoadInt32(&requests) != 1 || atomic.LoadInt32(&targetRequests) != 0 {
					t.Fatalf("redirect followed or status lost: %v", err)
				}
			})
		}
	}
}

func TestPaginationRefusesUnrelatedOriginsAndCycles(t *testing.T) {
	var targetRequests int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&targetRequests, 1)
	}))
	defer target.Close()
	for _, next := range []string{target.URL + "/transactions/", "//example.invalid/transactions/", "https://example.invalid/transactions/"} {
		t.Run(next, func(t *testing.T) {
			c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]interface{}{"next": next, "results": []interface{}{}})
			}, Options{})
			if _, err := c.ListAllTransactions(); err == nil || !strings.Contains(err.Error(), "origin") {
				t.Fatalf("unsafe next link accepted: %v", err)
			}
		})
	}
	if atomic.LoadInt32(&targetRequests) != 0 {
		t.Fatal("pagination leaked credentials to unrelated origin")
	}
	c := NewClient("token", "")
	for _, targetURL := range []string{
		"http://api.goshippo.com/transactions/",
		"https://api.goshippo.com:444/transactions/",
		"https://token@api.goshippo.com/transactions/",
		"https://api.goshippo.com.attacker.invalid/transactions/",
	} {
		if _, err := c.createRequest(context.Background(), "GET", targetURL, nil); err == nil {
			t.Fatalf("unsafe request URL accepted: %s", targetURL)
		}
	}
	var requests int32
	c = localClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		io.WriteString(w, "{\"next\":\"?results=25\",\"results\":[]}")
	}, Options{})
	if _, err := c.ListAllTransactions(); err == nil || !strings.Contains(err.Error(), "repeated") || requests != 1 {
		t.Fatalf("pagination cycle not stopped: %v (%d requests)", err, requests)
	}
}

func TestPaginationContextStopsBetweenPages(t *testing.T) {
	var requests int32
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		io.WriteString(w, "{\"next\":\"?page=2\",\"results\":[{}]}")
	}, Options{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := c.doListContext(ctx, "GET", "/transactions/", nil, func(json.RawMessage) error {
		cancel()
		return nil
	})
	if !stderrors.Is(err, context.Canceled) || requests != 1 {
		t.Fatalf("pagination ignored cancellation: %v / %d", err, requests)
	}
}

func TestResponseBodyTimeoutPreservesStatusAndCause(t *testing.T) {
	release := make(chan struct{})
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "partial")
		w.(http.Flusher).Flush()
		<-release
	}, Options{Timeout: 100 * time.Millisecond})
	defer close(release)
	_, err := c.RetrieveTransaction("id")
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.Status != 503 || string(apiErr.ResponseBody) != "partial" ||
		!stderrors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("lost status, partial response or timeout: %v", err)
	}
}

func TestTraceExcludesSecretsAndSupportsConcurrentClient(t *testing.T) {
	var logs bytes.Buffer
	logger := log.New(&logs, "", 0)
	c := localClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		io.WriteString(w, "secret-test-token private body")
	}, Options{})
	c.SetTraceLogger(logger)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.SetTraceLogger(logger)
			c.PurchaseShippingLabel(&models.TransactionInput{Rate: "private-rate", Metadata: "private-metadata"})
		}()
	}
	wg.Wait()
	output := logs.String()
	for _, secret := range []string{"secret-test-token", "Authorization", "private", "ShippoToken"} {
		if strings.Contains(output, secret) {
			t.Fatalf("trace contains %q", secret)
		}
	}
	if !strings.Contains(output, "status=400") || !strings.Contains(output, "method=POST") {
		t.Fatalf("trace lacks safe diagnostic information: %s", output)
	}
}
