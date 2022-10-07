package gohntest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetUpdates(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockItems := []int{1, 2, 3, 4, 5}
	mockProfiles := []string{"user1", "user2", "user3", "user4", "user5"}

	mockUpdates := gohn.Update{
		Items:    &mockItems,
		Profiles: &mockProfiles,
	}

	mockUpdatesJSON, err := json.Marshal(mockUpdates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mux.HandleFunc("/updates.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(mockUpdatesJSON))
	})

	ctx := context.Background()
	got, err := client.Updates.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error getting updates: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if got.Items == nil {
		t.Fatalf("expected items to be %v, got nil", mockItems)
	}

	if got.Profiles == nil {
		t.Fatalf("expected profiles to be %v, got nil", mockProfiles)
	}

	if len(*got.Items) != len(*mockUpdates.Items) {
		t.Errorf("expected updates %v, got %v", mockUpdates, *got)
	}

	if len(*got.Profiles) != len(*mockUpdates.Profiles) {
		t.Errorf("expected updates %v, got %v", mockUpdates, *got)
	}

	for i, id := range *got.Items {
		if id != (*mockUpdates.Items)[i] {
			t.Errorf("expected updates %v, got %v", mockUpdates, *got)
		}
	}

	for i, profile := range *got.Profiles {
		if profile != (*mockUpdates.Profiles)[i] {
			t.Errorf("expected updates %v, got %v", mockUpdates, *got)
		}
	}
}
