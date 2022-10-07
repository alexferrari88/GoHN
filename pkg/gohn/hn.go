package gohn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	BASE_URL   = "https://hacker-news.firebaseio.com/v0/"
	USER_AGENT = "GoHN/" + Version
)

// Client manages communication with the Hacker News API.
type Client struct {
	// HTTP client used to communicate with the API.
	httpClient *http.Client
	BaseURL    *url.URL

	UserAgent string

	common service

	// Services used for talking to different parts of the Hacker News API.
	Items   *ItemsService
	Stories *StoriesService
	Users   *UsersService
	Updates *UpdatesService
}

type service struct {
	client *Client
}

// GetHTTPClient returns the HTTP client used by the Client.
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// NewClient returns a new Client that will be used to make requests to the Hacker News API.
// If a nil httpClient is provided, http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(BASE_URL)
	c := Client{httpClient: httpClient, BaseURL: baseURL, UserAgent: USER_AGENT}
	c.common.client = &c
	c.Items = (*ItemsService)(&c.common)
	c.Stories = (*StoriesService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Updates = (*UpdatesService)(&c.common)
	return &c
}

// NewRequest creates an API request.
// path is a relative URL path (e.g. "items/1") and it will be resolved to the BaseURL of the Client.
func (c *Client) NewRequest(method, path string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL %q must have a trailing slash", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
// The Hacker News API returns JSON which is decoded and
// stored in the value pointed to by v, or returned as
// an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			// If the context was canceled, return the context's error,
			// which may be more useful than the underlying error.
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	// status codes between 200 and 299 are considered successful
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return &ErrResponse{Response: r}
}
