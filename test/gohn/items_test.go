package gohntest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
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

func generateBenchmarkItems(s, n int) []*gohn.Item {
	var items []*gohn.Item
	for i := s; i < n; i++ {
		i := i
		items = append(items, &gohn.Item{ID: &i})
	}
	return items
}

func BenchmarkFetchAllKids(b *testing.B) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	mockInitialKidsIDs := []int{2, 3, 4}
	mockExtraKids := generateBenchmarkItems(5, 10)
	mockKidID2KidsIDs := []int{5, 6}
	mockKidID2 := gohn.Item{ID: &mockInitialKidsIDs[0], Kids: &mockKidID2KidsIDs}
	mockKidID3KidsIDs := []int{7}
	mockKidID3 := gohn.Item{ID: &mockInitialKidsIDs[1], Kids: &mockKidID3KidsIDs}
	mockKidID4 := gohn.Item{ID: &mockInitialKidsIDs[2]}
	mockKids := []*gohn.Item{&mockKidID2, &mockKidID3, &mockKidID4}
	mockKids = append(mockKids, mockExtraKids[3:]...)

	var mockKidsIDs = new([]int)
	for _, kid := range mockKids {
		*mockKidsIDs = append(*mockKidsIDs, *kid.ID)
	}

	numDescendants := len(*mockKidsIDs) + 3

	mockType := "story"
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockType, Kids: mockKidsIDs, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": `+fmt.Sprintf("%d", numDescendants)+`}`)
	})

	// create mux.HandleFunc for each kid
	for _, kid := range mockKids {
		kidJSON, _ := json.Marshal(kid)
		mux.HandleFunc(fmt.Sprintf("/item/%v.json", *kid.ID), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, string(kidJSON))
		})
	}

	ctx := context.Background()
	for n := 0; n < b.N; n++ {
		_, _ = client.Items.FetchAllDescendants(ctx, mockParent, nil)
	}
}
