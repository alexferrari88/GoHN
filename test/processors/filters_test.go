package processorstest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestFilterOutWords_singleWord(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredWord := "potato"
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 7
	expectedNumDescendants := numDescendants - 1
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "`+filteredWord+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutWords([]string{filteredWord}, false))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.Text == nil {
			t.Error("expected text to be not nil")
		}
		if *item.Text == "potato" {
			t.Errorf("item with text '%s' should have been filtered out", filteredWord)
		}
	}
}

func TestFilterOutWords_singleWordButNotFound(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredWord := "potato"
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 7
	expectedNumDescendants := numDescendants
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutWords([]string{filteredWord}, false))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if *item.Text == "potato" {
			t.Errorf("item with text '%s' should have been filtered out", filteredWord)
		}
	}
}

func TestFilterOutWords_multipleWords(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredWords := []string{"potato", "tomato"}
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 8
	expectedNumDescendants := numDescendants - len(filteredWords)
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8, 9}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8, 9], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "`+filteredWords[0]+`"}`)
	})
	mux.HandleFunc("/item/9.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 9, "type": "comment", "text": "`+filteredWords[1]+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutWords(filteredWords, false))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if *item.Text == "potato" || *item.Text == "tomato" {
			t.Errorf("item with text '%s' or '%s' should have been filtered out", filteredWords[0], filteredWords[1])
		}
	}
}

func TestFilterOutWords_multipleWordsButOneNotFound(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredWords := []string{"potato", "tomato"}
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 7
	expectedNumDescendants := numDescendants - len(filteredWords) + 1
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "`+filteredWords[0]+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutWords(filteredWords, false))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if *item.Text == "potato" || *item.Text == "tomato" {
			t.Errorf("item with text '%s' or '%s' should have been filtered out", filteredWords[0], filteredWords[1])
		}
	}
}

func TestFilterOutWords_multipleWordsButAllNotFound(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredWords := []string{"potato", "tomato"}
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 6
	expectedNumDescendants := numDescendants
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutWords(filteredWords, false))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if *item.Text == "potato" || *item.Text == "tomato" {
			t.Errorf("item with text '%s' or '%s' should have been filtered out", filteredWords[0], filteredWords[1])
		}
	}
}

func TestFilterOutDeleted_multipleItems(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	mockParentType := "story"
	numDescendants := 8
	expectedNumDescendants := 6
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8, 9}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8, 9], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test", deleted: true}`)
	})
	mux.HandleFunc("/item/9.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 9, "type": "comment", "text": "test", deleted: true}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutDeleted())

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.Deleted != nil && *item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutDeleted_singleItem(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	mockParentType := "story"
	numDescendants := 8
	expectedNumDescendants := numDescendants - 1
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8, 9}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8, 9], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/9.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 9, "type": "comment", "text": "test", deleted: true}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutDeleted())

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.Deleted != nil && *item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutDeleted_noneFound(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	mockParentType := "story"
	numDescendants := 8
	expectedNumDescendants := numDescendants
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8, 9}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8, 9], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/9.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 9, "type": "comment", "text": "test"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutDeleted())

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.Deleted != nil && *item.Deleted {
			t.Errorf("item with deleted true should have been filtered out")
		}
	}
}

func TestFilterOutUsers_singleUser(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredUsersIds := []string{"bob"}
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 7
	expectedNumDescendants := numDescendants - len(filteredUsersIds)
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test", "by": "`+filteredUsersIds[0]+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutUsers(filteredUsersIds))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.By != nil && *item.By == filteredUsersIds[0] {
			t.Errorf("item with user '%s' should have been filtered out", filteredUsersIds[0])
		}
	}
}

func TestFilterOutUsers_multipleUsers(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	filteredUsersIds := []string{"bob", "ken"}
	mockParentID := 1
	mockParentType := "story"
	numDescendants := 7
	expectedNumDescendants := numDescendants - len(filteredUsersIds)
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}

	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 8], "descendants": `+fmt.Sprint(numDescendants)+`}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "text": "test"}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7], "text": "test"}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/5.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 5, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/6.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 6, "type": "comment", "text": "test"}`)
	})
	mux.HandleFunc("/item/7.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 7, "type": "comment", "text": "test", "by": "`+filteredUsersIds[0]+`"}`)
	})
	mux.HandleFunc("/item/8.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "test", "by": "`+filteredUsersIds[1]+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.FilterOutUsers(filteredUsersIds))

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != expectedNumDescendants {
		t.Errorf("expected %d items, got %v", expectedNumDescendants, len(got))
	}

	for _, item := range got {
		if item.By != nil && (*item.By == filteredUsersIds[0] || *item.By == filteredUsersIds[1]) {
			t.Errorf("item with user '%s' or '%s' should have been filtered out", filteredUsersIds[0], filteredUsersIds[1])
		}
	}
}
