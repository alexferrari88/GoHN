package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetTopStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetTopIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}

func TestGetBestStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/beststories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetBestIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}

func TestGetNewStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/newstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetNewIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}

func TestGetAskStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/askstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetAskIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}

func TestGetShowStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/showstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetShowIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}

func TestGetJobStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	mux.HandleFunc("/jobstories.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[1,2,3]")
	})

	mockStoriesIDs := []int{1, 2, 3}
	ctx := context.Background()

	got, err := client.Stories.GetJobIDs(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(got) != len(mockStoriesIDs) {
		t.Errorf("expected %v stories, got %v", len(mockStoriesIDs), len(got))
	}

	for i, story := range got {
		if *story != mockStoriesIDs[i] {
			t.Errorf("expected story %v, got %v", mockStoriesIDs[i], story)
		}
	}
}
