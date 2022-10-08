package gohntest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetMaxItemID(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 123

	mux.HandleFunc("/maxitem.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", mockID)
	})

	ctx := context.Background()
	got, err := client.Items.GetMaxID(ctx)

	if err != nil {
		t.Fatalf("unexpected error getting max item ID: %v", err)
	}

	if got == nil {
		t.Fatalf("expected max item ID to be %v, got nil", mockID)
	}

	if *got != mockID {
		t.Errorf("expected max item ID %d, got %d", mockID, got)
	}
}

func TestGetItem(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "title": "test title", "url": "http://example.com"}`)
	})

	ctx := context.Background()
	got, err := client.Items.Get(ctx, 1)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if *got.ID != 1 {
		t.Errorf("expected item ID %d, got %d", 1, got.ID)
	}
	if *got.Type != "story" {
		t.Errorf("expected item type %s, got %s", "story", *got.Type)
	}
	if *got.Title != "test title" {
		t.Errorf("expected item title %s, got %s", "test title", *got.Title)
	}
	if *got.URL != "http://example.com" {
		t.Errorf("expected item URL %s, got %s", "http://example.com", *got.URL)
	}
}

func TestFetchAllDescendants(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 1
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6]}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7]}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, nil)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 6 {
		t.Errorf("expected 6 items, got %v", len(got))
	}

	for _, id := range []int{2, 3, 4, 5, 6, 7} {
		if got[id] == nil {
			t.Fatalf("expected item %v to be %v, got nil", id, id)
		}
		if *got[id].ID != id {
			t.Errorf("expected item ID %d, got %d", id, *got[id].ID)
		}
	}
}

func TestSetCommentsPosition(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 1
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6]}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7]}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, nil)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 6 {
		t.Errorf("expected 6 items, got %v", len(got))
	}

	story := gohn.Story{
		Parent:          mockParent,
		CommentsByIdMap: got,
	}

	story.SetCommentsPosition()

	// sort story.CommentsIndex by the Position field
	var keys []int
	for k := range story.CommentsByIdMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return *(story.CommentsByIdMap[keys[i]].Position) < *(story.CommentsByIdMap[keys[j]].Position)
	})

	// expected order
	expectedPositionIDs := []int{2, 5, 6, 3, 7, 4}
	for i, id := range expectedPositionIDs {
		if id != keys[i] && *(story.CommentsByIdMap[id].Position) != i {
			t.Errorf("expected order %v, got %v", expectedPositionIDs, keys)
		}
	}
}

func TestIsTopLevelComment(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 1
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "parent": 1}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment",	"parent": 3}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, nil)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 6 {
		t.Errorf("expected 6 items, got %v", len(got))
	}

	story := gohn.Story{
		Parent:          mockParent,
		CommentsByIdMap: got,
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
		if isTop, _ := story.IsTopLevelComment(comment); isTop != expectedByID[*comment.ID] {
			t.Errorf("expected ID %d to be %v, got %v", *comment.ID, expectedByID[*comment.ID], isTop)
		}
	}
}

func TestGetOrderedCommentsIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 1
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "parent": 1}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment",	"parent": 3}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, nil)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 6 {
		t.Errorf("expected 6 items, got %v", len(got))
	}

	story := gohn.Story{
		Parent:          mockParent,
		CommentsByIdMap: got,
	}

	story.SetCommentsPosition()

	expectedIDs := []int{2, 5, 6, 3, 7, 4}
	orderedIDs, _ := story.GetOrderedCommentsIDs()

	if !reflect.DeepEqual(expectedIDs, orderedIDs) {
		t.Errorf("expected order %v, got %v", expectedIDs, orderedIDs)
	}
}

func TestGetStoryIdFromComment(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockID := 1
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "parent": 1}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment",	"parent": 3}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, nil)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 6 {
		t.Errorf("expected 6 items, got %v", len(got))
	}

	expectedStoryID := 1
	kidID := 6
	kidsParentId := 2
	kidsType := "comment"
	storyID, err := client.Items.GetStoryIdFromComment(ctx, &gohn.Item{ID: &kidID, Parent: &kidsParentId, Type: &kidsType})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if storyID == nil {
		t.Fatalf("expected storyID to be %v, got nil", expectedStoryID)
	}

	if expectedStoryID != *storyID {
		t.Errorf("expected story ID %v, got %v", expectedStoryID, storyID)
	}
}

func TestGetItem_nullResponse(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `null`)
	})

	ctx := context.Background()
	got, err := client.Items.Get(ctx, 1)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got != nil {
		t.Errorf("expected item to be nil, got %v", got)
	}
}

func TestFetchAllDescendants_Processor_excludeKids(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	idToExclude := 2
	excludedItemsIDs := []int{idToExclude, 5, 6}
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "parent": 1}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment",	"parent": 3}`)
	})

	ctx := context.Background()
	processor := func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		wg.Add(1)
		defer wg.Done()
		if item.ID != nil && *item.ID == idToExclude {
			return true, errors.New("mock error")
		}
		return false, nil
	}
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processor)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected items not to be nil")
	}

	if len(got) != len(excludedItemsIDs) {
		t.Errorf("expected %d items, got %v", len(excludedItemsIDs), len(got))
	}

	// check if any item with ID in excludedItemsIDs is in the result
	for _, item := range got {
		for _, excludedID := range excludedItemsIDs {
			if item.ID != nil && *item.ID == excludedID {
				t.Errorf("expected item with ID %v not to be in the result", excludedID)
			}
		}
	}
}

func TestFetchAllDescendants_Processor_includeKids(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	idToExclude := 2
	numDescendants := 6
	mockType := "story"
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 6}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "parent": 1}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "parent": 2}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment",	"parent": 3}`)
	})

	ctx := context.Background()
	processor := func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		wg.Add(1)
		defer wg.Done()
		if item.ID != nil && *item.ID == idToExclude {
			return false, errors.New("mock error")
		}
		return false, nil
	}
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processor)

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected items not to be nil")
	}

	if len(got) != numDescendants-1 {
		t.Errorf("expected %d items, got %v", numDescendants-1, len(got))
	}

	// check if the excluded item is in the result
	for _, item := range got {
		if item.ID != nil && *item.ID == idToExclude {
			t.Errorf("expected item with ID %v not to be in the result", idToExclude)
		}
	}
}
