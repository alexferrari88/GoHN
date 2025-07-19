package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/setup"
)

func TestGetStoriesIDs(t *testing.T) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	tests := []struct {
		storyType gohn.StoryType
		url       string
	}{
		{gohn.TopStory, "/topstories.json"},
		{gohn.BestStory, "/beststories.json"},
		{gohn.NewStory, "/newstories.json"},
		{gohn.AskStory, "/askstories.json"},
		{gohn.ShowStory, "/showstories.json"},
		{gohn.JobStory, "/jobstories.json"},
	}

	for _, tt := range tests {
		t.Run(string(tt.storyType), func(t *testing.T) {
			mux.HandleFunc(tt.url, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "[1,2,3]")
			})

			mockStoriesIDs := []int{1, 2, 3}
			ctx := context.Background()

			got, err := client.Stories.GetIDs(ctx, tt.storyType)
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
		})
	}
}