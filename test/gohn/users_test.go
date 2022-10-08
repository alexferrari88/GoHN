package gohntest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetUser(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockUserId := "testuser"
	mockUser := struct {
		ID string `json:"id"`
	}{
		ID: mockUserId,
	}

	mockUserJSON, err := json.Marshal(mockUser)

	if err != nil {
		t.Fatalf("error marshalling mock user: %v", err)
	}

	mux.HandleFunc(fmt.Sprintf("/user/%s.json", mockUserId), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(mockUserJSON))
	})

	ctx := context.Background()
	got, err := client.Users.GetByUsername(ctx, mockUserId)

	if err != nil {
		t.Fatalf("unexpected error getting user: %v", err)
	}

	if got == nil {
		t.Fatalf("expected user to be %v, got nil", mockUserId)
	}

	if *got.ID != mockUserId {
		t.Errorf("expected user id %s, got %v", mockUserId, *got.ID)
	}
}
