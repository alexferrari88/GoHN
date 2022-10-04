package processorstest

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
	"github.com/alexferrari88/gohn/test/mocks"
)

func setupItems(extraKids []gohn.Item) (gohn.Item, []gohn.Item) {
	mockItem := gohn.Item{
		ID:   1,
		Kids: []int{2, 3, 4},
	}
	mockKid1 := gohn.Item{
		ID:   2,
		Text: "test",
		Kids: []int{5, 6},
	}
	mockKid2 := gohn.Item{
		ID:   3,
		Text: "test",
		Kids: []int{7},
	}
	mockKid3 := gohn.Item{
		ID:   4,
		Text: "test",
	}
	mockKid4 := gohn.Item{
		ID:   5,
		Text: "test",
	}
	mockKid5 := gohn.Item{
		ID:   6,
		Text: "test",
	}
	mockKid6 := gohn.Item{
		ID:   7,
		Text: "test",
	}

	kids := []gohn.Item{
		mockKid1,
		mockKid2,
		mockKid3,
		mockKid4,
		mockKid5,
		mockKid6,
	}

	if len(extraKids) > 0 {
		for _, kid := range extraKids {
			mockItem.Kids = append(mockItem.Kids, kid.ID)
			kids = append(kids, kid)
		}
	}

	return mockItem, kids

}

func setup(extraKids []gohn.Item) (gohn.Item, []gohn.Item, *mocks.MockHTTPClient, error) {
	mockItem, kids := setupItems(extraKids)
	mockClient, err := mocks.SetupMockClient(mockItem, kids)
	if err != nil {
		return gohn.Item{}, nil, nil, fmt.Errorf("error setting up test: %v", err)
	}
	return mockItem, kids, mockClient, nil
}

func TestFilterOutWordsSingleWord(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:   8,
			Text: "potato",
		},
	})

	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredWord := "potato"
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutWords([]string{filteredWord}, false))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Text == "potato" {
			t.Errorf("item with text 'potato' should have been filtered out")
		}
	}
}

func TestFilterOutWordsSingleWordButNotFound(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{})

	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredWord := "potato"
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutWords([]string{filteredWord}, false))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Text == "potato" {
			t.Errorf("item with text 'potato' should have been filtered out")
		}
	}
}

func TestFilterOutWordsMultipleWords(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:   8,
			Text: "potato",
		},
		{
			ID:   9,
			Text: "tomato",
		},
	})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredWords := []string{"potato", "tomato"}
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutWords(filteredWords, false))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Text == "potato" || item.Text == "tomato" {
			t.Errorf("item with text 'potato' or 'tomato' should have been filtered out")
		}
	}
}

func TestFilterOutWordsMultipleWordsButOneNotFound(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:   8,
			Text: "potato",
		}})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredWords := []string{"potato", "tomato"}
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutWords(filteredWords, false))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Text == "potato" || item.Text == "tomato" {
			t.Errorf("item with text 'potato' or 'tomato' should have been filtered out")
		}
	}
}

func TestFilterOutWordsMultipleWordsButAllNotFound(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredWords := []string{"potato", "tomato"}
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutWords(filteredWords, false))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Text == "potato" || item.Text == "tomato" {
			t.Errorf("item with text 'potato' or 'tomato' should have been filtered out")
		}
	}
}

func TestFilterOutDeletedMultipleItems(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:      8,
			Text:    "test",
			Deleted: true,
		},
		{
			ID:      9,
			Text:    "test",
			Deleted: true,
		},
	})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutDeleted())

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutDeletedSingleItem(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:      8,
			Text:    "test",
			Deleted: true,
		},
	})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutDeleted())

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutDeletedButNoneFound(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutDeleted())

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutUsersSingleUser(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID: 8,
			By: "test",
		},
	})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredUsers := []string{"test"}
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutUsers(filteredUsers))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.By == "test" {
			t.Errorf("item with user 'test' should have been filtered out")
		}
	}
}

func TestFilterOutUsersMultipleUsers(t *testing.T) {
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID: 8,
			By: "test",
		},
		{
			ID: 9,
			By: "test2",
		},
	})
	if err != nil {
		t.Errorf("error setting up test: %v", err)
	}

	filteredUsers := []string{"test", "test2"}
	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.FilterOutUsers(filteredUsers))

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, item := range items {
		if item.By == "test" || item.By == "test2" {
			t.Errorf("item with user 'test' or 'test2' should have been filtered out")
		}
	}
}
