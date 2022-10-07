package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetUser(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockUserId := "testuser"

	mux.HandleFunc(fmt.Sprintf("/user/%s.json", mockUserId), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf(`{"id": "%s"}`, mockUserId))
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
