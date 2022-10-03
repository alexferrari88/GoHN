package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/mocks"
)

func TestGetUser(t *testing.T) {
	mockUser := gohn.User{
		ID: "test",
	}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockUser)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{fmt.Sprintf(gohn.USER_URL, mockUser.ID)}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	user, err := client.GetUser("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.ID != mockUser.ID {
		t.Errorf("expected user %v, got %v", mockUser, user)
	}
}
