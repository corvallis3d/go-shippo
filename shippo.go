package shippo

import "github.com/corvallis3d/go-shippo/client"

const (
	APIVersion20140211 = "2014-02-11"
	// Deprecated: Use APIVersion20161025.
	APIVesrion20161025 = "2016-10-25"
	APIVersion20161025 = APIVesrion20161025
	APIVersion20170329 = "2017-03-29"
	APIVersion20180208 = "2018-02-08"
)

// NewClient creates a new Shippo client using the account's default API version.
// Requests have client.DefaultTimeout; use NewClientWithOptions to configure it.
func NewClient(privateToken string) *client.Client {
	return client.NewClient(privateToken, "")
}

// NewClientWithVersion creates a new Shippo client with API version explicitly specified.
func NewClientWithVersion(privateToken, apiVersion string) *client.Client {
	return client.NewClient(privateToken, apiVersion)
}

// NewClientWithOptions creates a reusable client with HTTP configuration.
// Pass an explicit API version, or an empty string to use the account default.
func NewClientWithOptions(privateToken, apiVersion string, options client.Options) *client.Client {
	return client.NewClientWithOptions(privateToken, apiVersion, options)
}
