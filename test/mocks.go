package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MockHTTPClient struct {
	Responses map[string]*http.Response
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if resp, ok := m.Responses[req.URL.String()]; ok {
		return resp, nil
	} else {
		return nil, fmt.Errorf("no response found for request: %v", req)
	}
}

func NewMockClient(urls []string, responses []*http.Response) *MockHTTPClient {
	m := &MockHTTPClient{
		Responses: make(map[string]*http.Response),
	}

	for i, url := range urls {
		m.Responses[url] = responses[i]
	}

	return m
}

func NewMockResponse(statusCode int, body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}, nil
}
