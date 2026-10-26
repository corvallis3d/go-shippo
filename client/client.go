package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/corvallis3d/go-shippo/errors"
	"github.com/corvallis3d/go-shippo/models"
)

const shippoAPIBaseURL = "https://api.goshippo.com"

// DefaultTimeout bounds each HTTP request, including reading the response body.
const DefaultTimeout = 30 * time.Second

// Options configures a reusable client. Configure the transport before use;
// custom transports must honor request contexts and must not retry writes.
type Options struct {
	// HTTPClient is copied, retaining its transport and connection pool. Redirects
	// are always disabled to prevent credential forwarding and write replay.
	HTTPClient *http.Client
	// Timeout overrides HTTPClient.Timeout when positive. Otherwise its positive
	// timeout is retained, or DefaultTimeout is used. Contexts can end requests sooner.
	Timeout time.Duration
}

// Client can be shared by concurrent callers.
type Client struct {
	privateToken string
	apiVersion   string
	httpClient   *http.Client
	baseURL      string
	loggerMu     sync.RWMutex
	logger       *log.Logger
}

type listOutputCallback func(v json.RawMessage) error

// NewClient creates a client with a bounded default timeout. An empty apiVersion
// leaves Shippo's account default in effect, as in earlier releases.
func NewClient(privateToken, apiVersion string) *Client {
	return NewClientWithOptions(privateToken, apiVersion, Options{})
}

// NewClientWithOptions creates a client with an explicit API version and HTTP
// configuration. It does not modify the supplied http.Client or its transport.
func NewClientWithOptions(privateToken, apiVersion string, options Options) *Client {
	httpClient := &http.Client{}
	if options.HTTPClient != nil {
		*httpClient = *options.HTTPClient
	}
	if options.Timeout > 0 {
		httpClient.Timeout = options.Timeout
	} else if httpClient.Timeout <= 0 {
		httpClient.Timeout = DefaultTimeout
	}
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{
		privateToken: privateToken,
		apiVersion:   apiVersion,
		httpClient:   httpClient,
		baseURL:      shippoAPIBaseURL,
	}
}

// SetTraceLogger sets a trace logger and returns the previous logger. Only HTTP
// methods and status codes are logged; credentials, URLs and bodies are excluded.
func (c *Client) SetTraceLogger(logger *log.Logger) *log.Logger {
	c.loggerMu.Lock()
	defer c.loggerMu.Unlock()
	oldLogger := c.logger
	c.logger = logger
	return oldLogger
}

func (c *Client) doContext(ctx context.Context, method, path string, input, output interface{}) error {
	req, err := c.createRequest(ctx, method, c.baseURL+path, input)
	if err != nil {
		return err
	}
	return c.executeRequest(req, output)
}

// doListContext is shared by all full-list helpers. The caller's context bounds
// the entire traversal; the HTTP timeout applies separately to each page.
func (c *Client) doListContext(ctx context.Context, method, path string, input interface{}, outputCallback listOutputCallback) error {
	nextURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("parse list URL: %w", err)
	}
	filters := nextURL.Query()
	if !filters.Has("results") {
		filters.Set("results", "25")
	}
	nextURL.RawQuery = filters.Encode()
	visited := make(map[string]bool)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if visited[nextURL.String()] {
			return fmt.Errorf("pagination returned a repeated page")
		}
		visited[nextURL.String()] = true
		req, err := c.createRequest(ctx, method, nextURL.String(), input)
		if err != nil {
			return err
		}
		listOutput := &models.ListAPIOutput{}
		if err := c.executeRequest(req, listOutput); err != nil {
			return err
		}
		for _, v := range listOutput.Results {
			if err := outputCallback(v); err != nil {
				return fmt.Errorf("unmarshal list item: %w", err)
			}
		}
		if listOutput.NextPageURL == nil || *listOutput.NextPageURL == "" {
			return nil
		}
		next, err := url.Parse(*listOutput.NextPageURL)
		if err != nil {
			return fmt.Errorf("parse next page URL: %w", err)
		}
		nextURL = nextURL.ResolveReference(next)
		query := nextURL.Query()
		for key, values := range filters {
			if key != "page" && key != "cursor" && !query.Has(key) {
				query[key] = values
			}
		}
		nextURL.RawQuery = query.Encode()
	}
}

func (c *Client) createRequest(ctx context.Context, method, requestURL string, bodyObject interface{}) (*http.Request, error) {
	target, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("parse request URL: %w", err)
	}
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse API URL: %w", err)
	}
	if target.User != nil || !strings.EqualFold(target.Scheme, base.Scheme) || !strings.EqualFold(target.Host, base.Host) {
		return nil, fmt.Errorf("refusing authenticated request outside the API origin")
	}
	var body io.Reader
	if bodyObject != nil {
		data, err := json.Marshal(bodyObject)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, target.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}
	// Do not make write bodies replayable by net/http. Purchases are attempted once.
	req.GetBody = nil
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "ShippoToken "+c.privateToken)
	if c.apiVersion != "" {
		req.Header.Set("Shippo-API-Version", c.apiVersion)
	}
	return req, nil
}

func (c *Client) executeRequest(req *http.Request, output interface{}) error {
	c.logPrintf("Shippo HTTP request: method=%s", req.Method)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("make HTTP request: %w", err)
	}
	defer res.Body.Close()
	c.logPrintf("Shippo HTTP response: status=%d", res.StatusCode)
	data, readErr := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return &errors.APIError{Status: res.StatusCode, ResponseBody: data, Err: readErr}
	}
	if readErr != nil {
		return fmt.Errorf("read HTTP %d response: %w", res.StatusCode, readErr)
	}
	if output != nil && len(data) > 0 {
		if err := json.Unmarshal(data, output); err != nil {
			return fmt.Errorf("unmarshal HTTP %d response: %w", res.StatusCode, err)
		}
	}
	return nil
}

func (c *Client) logPrintf(format string, args ...interface{}) {
	c.loggerMu.RLock()
	logger := c.logger
	c.loggerMu.RUnlock()
	if logger != nil {
		logger.Printf(format, args...)
	}
}
