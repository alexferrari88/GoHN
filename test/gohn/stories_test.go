package gohntest

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/mocks"
)

func TestGetTopStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.TOP_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetTopStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}

func TestGetBestStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.BEST_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetBestStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}

func TestGetNewStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.NEW_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetNewStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}

func TestGetAskStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.ASK_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetAskStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}

func TestGetShowStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.SHOW_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetShowStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}

func TestGetJobStoriesIDs(t *testing.T) {
	mockStories := []int{1, 2, 3}
	mockResponseJSON, err := mocks.NewMockResponse(http.StatusOK, mockStories)
	if err != nil {
		t.Errorf("error creating mock response: %v", err)
	}
	mockClient := mocks.NewMockClient([]string{gohn.JOB_STORIES_URL}, []*http.Response{mockResponseJSON})

	client := gohn.NewClient(context.Background(), mockClient)
	stories, err := client.GetJobStoriesIDs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stories) != len(mockStories) {
		t.Errorf("expected %v stories, got %v", len(mockStories), len(stories))
	}

	for i, story := range stories {
		if story != mockStories[i] {
			t.Errorf("expected story %v, got %v", mockStories[i], story)
		}
	}
}
