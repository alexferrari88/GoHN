package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestNewClient(t *testing.T) {
	c := gohn.NewClient(nil)

	if c.BaseURL.String() != gohn.BASE_URL {
		t.Errorf("expected base url %v, got %v", gohn.BASE_URL, c.BaseURL)
	}

	if c.UserAgent != gohn.USER_AGENT {
		t.Errorf("expected user agent %v, got %v", gohn.USER_AGENT, c.UserAgent)
	}

	c2 := gohn.NewClient(nil)
	if c.GetHTTPClient() == c2.GetHTTPClient() {
		t.Errorf("expected different http.Clients, got the same")
	}
}

func TestNewRequest(t *testing.T) {
	c := gohn.NewClient(nil)

	inURL, outURL := "test", gohn.BASE_URL+"test"
	req, _ := c.NewRequest("GET", inURL)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("expected request URL %v, got %v", want, got)
	}

	userAgent := req.Header.Get("User-Agent")

	// test that default user-agent is attached to the request
	if got, want := userAgent, c.UserAgent; got != want {
		t.Errorf("expected user agent %v, got %v", want, got)
	}

	if !strings.Contains(userAgent, gohn.Version) {
		t.Errorf("User-Agent should contain %v, found %v", gohn.Version, userAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := gohn.NewClient(nil)
	_, err := c.NewRequest("GET", ":")
	if err == nil {
		t.Error("expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}
func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api", wantError: true},
		{rawurl: "https://example.com/api/", wantError: false},
	}
	c := gohn.NewClient(nil)
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewRequest(http.MethodGet, "test"); test.wantError && err == nil {
			t.Fatalf("expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("expected error to be nil, got %v.", err)
		}
	}
}

func TestDo(t *testing.T) {
	c, mux, _, teardown := setup.Init()
	defer teardown()

	type foo struct {
		Foo string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"Foo": "bar"}`)
	})

	req, _ := c.NewRequest("GET", ".")

	body := new(foo)
	ctx := context.Background()
	_, err := c.Do(ctx, req, body)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if body.Foo != "bar" {
		t.Errorf("expected body.Foo to be bar, got %v", body.Foo)
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	req, _ := client.NewRequest("GET", ".")
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatalf("Expected HTTP %d error, got no error.", http.StatusForbidden)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected HTTP %d error, got %d status code.", http.StatusForbidden, resp.StatusCode)
	}
}

func TestDo_bodyNull(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `null`)
	})

	req, _ := client.NewRequest("GET", ".")
	ctx := context.Background()
	resp, err := client.Do(ctx, req, nil)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP %d, got %d status code.", http.StatusOK, resp.StatusCode)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
	}
	err := gohn.CheckResponse(res).(*gohn.ErrResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &gohn.ErrResponse{
		Response: res,
	}
	if err == want {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}
