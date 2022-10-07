package processorstest

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestUnescapeHTML(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mockParentID := 1
	numDescendants := 7
	mockParentType := "story"
	mockParent := &gohn.Item{ID: &mockParentID, Type: &mockParentType, Kids: &[]int{2, 3, 4, 8}, Descendants: &numDescendants}
	mockItemToTestID := 8
	expectedText := `This is an <a href="https://www.example.com">example</a>`

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
		fmt.Fprint(w, `{"id": 8, "type": "comment", "text": "`+html.EscapeString(expectedText)+`"}`)
	})

	ctx := context.Background()
	got, err := client.Items.FetchAllDescendants(ctx, mockParent, processors.UnescapeHTML())

	if err != nil {
		t.Fatalf("unexpected error getting item: %v", err)
	}

	if got == nil {
		t.Fatalf("expected item to be %v, got nil", 1)
	}

	if len(got) != 7 {
		t.Fatalf("expected 7 items, got %v", len(got))
	}

	if *got[mockItemToTestID].Text != expectedText {
		t.Errorf("expected unescaped text, got %v", *got[mockItemToTestID].Text)
	}

}
