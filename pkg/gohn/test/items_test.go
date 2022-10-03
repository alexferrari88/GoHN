package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

func TestGetMaxItemID(t *testing.T) {
	mockID := 123
	mockResponseJSON, err := NewMockResponse(http.StatusOK, mockID)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := NewMockClient([]string{gohn.MAX_ITEM_ID_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	id, err := client.GetMaxItemID()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if id != mockID {
		t.Errorf("expected id %v, got %v", mockID, id)
	}
}

func TestGetItem(t *testing.T) {
	mockItem := gohn.Item{
		ID: 1,
	}
	mockResponseJSON, err := NewMockResponse(http.StatusOK, mockItem)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := NewMockClient([]string{fmt.Sprintf(gohn.ITEM_URL, 1)}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	item, err := client.GetItem(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if item.ID != mockItem.ID {
		t.Errorf("expected item %v, got %v", mockItem, item)
	}
}

// test the RetrieveKidsItems function
// item has 3 kids
// kid 1 has 2 kids
// kid 2 has 1 kid
// kid 3 has 0 kids
// total number of items should be 6
func TestRetrieveKidsItems(t *testing.T) {
	mockItem := gohn.Item{
		ID:   1,
		Kids: []int{2, 3, 4},
	}
	mockKid1 := gohn.Item{
		ID:   2,
		Kids: []int{5, 6},
	}
	mockKid2 := gohn.Item{
		ID:   3,
		Kids: []int{7},
	}
	mockKid3 := gohn.Item{
		ID: 4,
	}
	mockKid4 := gohn.Item{
		ID: 5,
	}
	mockKid5 := gohn.Item{
		ID: 6,
	}
	mockKid6 := gohn.Item{
		ID: 7,
	}
	mockResponseJSON, err := NewMockResponse(http.StatusOK, mockItem)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON2, err := NewMockResponse(http.StatusOK, mockKid1)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON3, err := NewMockResponse(http.StatusOK, mockKid2)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON4, err := NewMockResponse(http.StatusOK, mockKid3)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON5, err := NewMockResponse(http.StatusOK, mockKid4)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON6, err := NewMockResponse(http.StatusOK, mockKid5)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockResponseJSON7, err := NewMockResponse(http.StatusOK, mockKid6)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := NewMockClient([]string{
		fmt.Sprintf(gohn.ITEM_URL, 1),
		fmt.Sprintf(gohn.ITEM_URL, 2),
		fmt.Sprintf(gohn.ITEM_URL, 3),
		fmt.Sprintf(gohn.ITEM_URL, 4),
		fmt.Sprintf(gohn.ITEM_URL, 5),
		fmt.Sprintf(gohn.ITEM_URL, 6),
		fmt.Sprintf(gohn.ITEM_URL, 7),
	}, []*http.Response{
		mockResponseJSON,
		mockResponseJSON2,
		mockResponseJSON3,
		mockResponseJSON4,
		mockResponseJSON5,
		mockResponseJSON6,
		mockResponseJSON7,
	})

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, nil)

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	if items[mockKid1.ID].ID != mockKid1.ID {
		t.Errorf("expected item %v, got %v", mockKid1, items[mockKid1.ID])
	}

	if items[mockKid2.ID].ID != mockKid2.ID {
		t.Errorf("expected item %v, got %v", mockKid2, items[mockKid2.ID])
	}

	if items[mockKid3.ID].ID != mockKid3.ID {
		t.Errorf("expected item %v, got %v", mockKid3, items[mockKid3.ID])
	}

	if items[mockKid4.ID].ID != mockKid4.ID {
		t.Errorf("expected item %v, got %v", mockKid4, items[mockKid4.ID])
	}

	if items[mockKid5.ID].ID != mockKid5.ID {
		t.Errorf("expected item %v, got %v", mockKid5, items[mockKid5.ID])
	}

	if items[mockKid6.ID].ID != mockKid6.ID {
		t.Errorf("expected item %v, got %v", mockKid6, items[mockKid6.ID])
	}
}
