package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/mocks"
)

func setupMocks() (parent gohn.Item, kids []gohn.Item, mockClient *mocks.MockHTTPClient, err error) {
	// parent is the main item with ID = 1. It has 3 kids
	// kid[0] (ID = 2) has 2 kids (kid[3] (ID = 5) and kid[4] (ID = 6))
	// kid[1] (ID = 3) has 1 kid (kid[5] (ID = 7))
	// kid[2] (ID = 4) has 0 kids
	mockItems := mocks.NewMockItems(7)
	parent = mockItems[0]
	kids = mockItems[1:]
	mocks.AddKidsToMockItem(&parent, kids[0:3])
	mocks.AddKidsToMockItem(&kids[0], kids[3:5])
	mocks.AddKidsToMockItem(&kids[1], kids[5:6])

	mockClient, err = mocks.SetupMockClient(parent, kids)

	if err != nil {
		return parent, kids, mockClient, fmt.Errorf("error setting up mock client: %v", err)
	}
	return
}

func TestGetMaxItemID(t *testing.T) {
	mockID := 123
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockID)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.MAX_ITEM_ID_URL}, []*http.Response{mockResponseJSON})

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
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockItem)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{fmt.Sprintf(gohn.ITEM_URL, 1)}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	item, err := client.GetItem(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if item.ID != mockItem.ID {
		t.Errorf("expected item %v, got %v", mockItem, item)
	}
}

func TestRetrieveKidsItems(t *testing.T) {
	parent, kids, mockClient, err := setupMocks()

	if err != nil {
		t.Error(err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(parent, nil)

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %v", len(items))
	}

	for _, kid := range kids {
		if _, ok := items[kid.ID]; !ok {
			t.Errorf("expected item %v, got %v", kid, items)
		}
	}
}

func TestCalculateCommentsPosition(t *testing.T) {
	parent, _, mockClient, err := setupMocks()

	if err != nil {
		t.Error(err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(parent, nil)

	if len(items) != 6 {
		t.Fatalf("expected 6 items, got %v", len(items))
	}

	story := gohn.Story{
		StoryItem:       parent,
		CommentsByIdMap: items,
	}

	story.SetCommentsPosition()

	// sort story.CommentsIndex by the Position field
	var keys []int
	for k := range story.CommentsByIdMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return story.CommentsByIdMap[keys[i]].Position < story.CommentsByIdMap[keys[j]].Position
	})

	// expected order
	expectedPositionIDs := []int{2, 5, 6, 3, 7, 4}
	for i, id := range expectedPositionIDs {
		if id != keys[i] && story.CommentsByIdMap[id].Position != i {
			t.Errorf("expected order %v, got %v", expectedPositionIDs, keys)
		}
	}
}

func TestIsTopLevelComment(t *testing.T) {
	parent, _, mockClient, err := setupMocks()

	if err != nil {
		t.Error(err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(parent, nil)

	if len(items) != 6 {
		t.Fatalf("expected 6 items, got %v", len(items))
	}

	story := gohn.Story{
		StoryItem:       parent,
		CommentsByIdMap: items,
	}

	expectedByID := map[int]bool{
		2: true,
		3: true,
		4: true,
		5: false,
		6: false,
		7: false,
	}

	for _, comment := range story.CommentsByIdMap {
		if story.IsTopLevelComment(comment) != expectedByID[comment.ID] {
			t.Errorf("expected ID %d to be %v, got %v", comment.ID, expectedByID[comment.ID], story.IsTopLevelComment(comment))
		}
	}
}

func TestGetOrderedCommentsIDs(t *testing.T) {
	parent, _, mockClient, err := setupMocks()

	if err != nil {
		t.Error(err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(parent, nil)

	if len(items) != 6 {
		t.Fatalf("expected 6 items, got %v", len(items))
	}

	story := gohn.Story{
		StoryItem:       parent,
		CommentsByIdMap: items,
	}

	story.SetCommentsPosition()

	expectedIDs := []int{2, 5, 6, 3, 7, 4}
	orderedIDs := story.GetOrderedCommentsIDs()

	if !reflect.DeepEqual(expectedIDs, orderedIDs) {
		t.Errorf("expected order %v, got %v", expectedIDs, orderedIDs)
	}
}
