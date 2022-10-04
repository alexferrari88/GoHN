package mocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alexferrari88/gohn/pkg/gohn"
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

func NewMockItems(num int) []gohn.Item {
	items := make([]gohn.Item, num)
	for i := 0; i < num; i++ {
		items[i] = gohn.Item{
			ID:   i + 1,
			Text: "test",
		}
	}
	return items
}

func AddParentToMockItem(item *gohn.Item, parent *gohn.Item) {
	item.Parent = parent.ID
}

func AddKidsToMockItem(item *gohn.Item, kids []gohn.Item) {
	for _, kid := range kids {
		kid := kid
		AddParentToMockItem(&kid, item)
		item.Kids = append(item.Kids, kid.ID)
	}
}

func SetupMockClient(mockItem gohn.Item, mockKids []gohn.Item) (*MockHTTPClient, error) {

	mockResponseJSON, err := NewMockResponse(http.StatusOK, mockItem)
	if err != nil {
		return nil, fmt.Errorf("error creating mock response: %v", err)
	}

	mockKidResponses := make([]*http.Response, len(mockKids))
	for i, kid := range mockKids {
		mockKidResponses[i], err = NewMockResponse(http.StatusOK, kid)
		if err != nil {
			return nil, fmt.Errorf("error creating mock response: %v", err)
		}
	}

	urls := make([]string, len(mockKids)+1)
	urls[0] = fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", mockItem.ID)
	for i, kid := range mockKids {
		urls[i+1] = fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", kid.ID)
	}

	mockClient := NewMockClient(urls, append([]*http.Response{mockResponseJSON}, mockKidResponses...))

	return mockClient, nil
}
