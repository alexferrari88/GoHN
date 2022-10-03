package gohn

import (
	"context"
	"io/ioutil"
	"net/http"
)

// HTTPClient is an interface that is satisfied by http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type client struct {
	ctx    context.Context
	client HTTPClient
}

// NewClient returns a new Client that will be used to make requests to the Hacker News API.
// It uses context.Background() as the default context.
func NewClient(ctx context.Context, httpClient HTTPClient) *client {
	return &client{ctx: context.Background(), client: httpClient}
}

// NewClientWithContext returns a new Client that will be used to make requests to the Hacker News API
// It uses the provided context for HTTP requests.
func NewClientWithContext(ctx context.Context, httpClient HTTPClient) *client {
	return &client{ctx: ctx, client: httpClient}
}

// NewDefaultClient returns a http.DefaultClient that will be used to make requests to the Hacker News API.
// If a nil ctx is provided, context.Background() will be used for HTTP requests.
func NewDefaultClient() *client {
	return &client{ctx: context.Background(), client: http.DefaultClient}
}

// NewDefaultClientWithContext returns a new client to make requests to the Hacker News API.
// It uses a http.DefaultClient and the provided context for HTTP requests.
func NewDefaultClientWithContext(ctx context.Context) *client {
	return &client{ctx: ctx, client: http.DefaultClient}
}

// SetContext sets the context for HTTP requests.
func (c *client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// retrieveFromURL sends a GET request to the given URL
// and returns the raw response body values and an error.
func (c client) retrieveFromURL(url string) ([]byte, error) {
	var body []byte

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
