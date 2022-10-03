package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

func TestGetUpdates(t *testing.T) {
	mockUpdates := gohn.Update{
		Items:    []int{1, 2, 3},
		Profiles: []string{"user1", "user2"},
	}
	mockResponseJSON, err := NewMockResponse(http.StatusOK, mockUpdates)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := NewMockClient([]string{gohn.UPDATES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	updates, err := client.GetUpdates()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(updates.Items) != len(mockUpdates.Items) {
		t.Errorf("expected updates %v, got %v", mockUpdates, updates)
	}

	if len(updates.Profiles) != len(mockUpdates.Profiles) {
		t.Errorf("expected updates %v, got %v", mockUpdates, updates)
	}

	for i, id := range updates.Items {
		if id != mockUpdates.Items[i] {
			t.Errorf("expected updates %v, got %v", mockUpdates, updates)
		}
	}

	for i, profile := range updates.Profiles {
		if profile != mockUpdates.Profiles[i] {
			t.Errorf("expected updates %v, got %v", mockUpdates, updates)
		}
	}
}
