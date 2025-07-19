package gohntest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/test/setup"
)

// BenchmarkFetchAllDescendants_HNThread34591291 benchmarks fetching all comments
// from HN thread 34591291 which has 271 comments.
// This benchmark tests the real-world performance of the comment fetching system.
func BenchmarkFetchAllDescendants_HNThread34591291(b *testing.B) {
	client, err := gohn.NewClient(&http.Client{
		Timeout: time.Second * 30,
	})
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	
	// Get the actual HN thread item
	story, err := client.Items.Get(ctx, 34591291)
	if err != nil {
		b.Fatalf("Failed to get HN thread 34591291: %v", err)
	}

	if story.Kids == nil {
		b.Skip("Story has no comments to fetch")
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		comments, err := client.Items.FetchAllDescendants(ctx, story, nil)
		if err != nil {
			b.Fatalf("Failed to fetch comments: %v", err)
		}
		_ = comments // Prevent optimization
	}
}

// BenchmarkFetchAllDescendants_Mock271Comments creates a mock server with 271 comments
// to test performance without rate limiting from the real HN API.
func BenchmarkFetchAllDescendants_Mock271Comments(b *testing.B) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	// Setup mock responses for a story with 271 comments in a tree structure
	// This simulates the HN thread structure but without real API rate limits
	setupMock271Comments(mux)

	mockStory := &gohn.Item{
		ID:          intPtr(1),
		Type:        stringPtr("story"),
		Kids:        &[]int{2, 3, 4, 5, 6, 7, 8, 9, 10}, // Top-level comments
		Descendants: intPtr(271),
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		comments, err := client.Items.FetchAllDescendants(ctx, mockStory, nil)
		if err != nil {
			b.Fatalf("Failed to fetch comments: %v", err)
		}
		
		// Verify we got the expected number of comments
		if len(comments) != 271 {
			b.Fatalf("Expected 271 comments, got %d", len(comments))
		}
	}
}

// setupMock271Comments creates a mock server setup that simulates
// a HN thread with 271 comments in a realistic tree structure
func setupMock271Comments(mux *http.ServeMux) {
	// Root story
	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4, 5, 6, 7, 8, 9, 10], "descendants": 271}`)
	})

	// Generate a realistic comment tree structure
	commentID := 2
	
	// Top-level comments (9 comments, IDs 2-10)
	for i := 2; i <= 10; i++ {
		kids := generateKidsForComment(i, &commentID, 30) // Each top-level has ~30 descendants
		kidsJSON := formatKidsArray(kids)
		localID := i // Capture for closure
		localKidsJSON := kidsJSON // Capture for closure
		
		mux.HandleFunc(fmt.Sprintf("/item/%d.json", localID), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"id": %d, "type": "comment", "parent": 1, "kids": %s}`, localID, localKidsJSON)
		})
	}

	// Generate all descendant comments
	for id := 11; id <= 271; id++ {
		// Randomly decide if this comment has children (about 30% chance)
		var kidsJSON string = "null"
		if id%3 == 0 && id < 250 { // Don't let the last 21 comments have kids to control total
			if id+200 <= 271 {
				kidsJSON = fmt.Sprintf("[%d, %d]", id+100, id+200)
			} else if id+100 <= 271 {
				kidsJSON = fmt.Sprintf("[%d]", id+100)
			}
		}
		
		parentID := calculateParentID(id)
		localID := id // Capture for closure
		localParentID := parentID // Capture for closure
		localKidsJSON := kidsJSON // Capture for closure
		
		mux.HandleFunc(fmt.Sprintf("/item/%d.json", localID), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"id": %d, "type": "comment", "parent": %d, "kids": %s}`, localID, localParentID, localKidsJSON)
		})
	}
}

// generateKidsForComment generates child comment IDs for a given comment
func generateKidsForComment(parentID int, currentID *int, maxKids int) []int {
	numKids := min(maxKids, 5) // Limit direct children to 5
	kids := make([]int, numKids)
	
	for i := 0; i < numKids; i++ {
		*currentID++
		kids[i] = *currentID
	}
	
	return kids
}

// calculateParentID calculates the parent ID for a given comment ID
func calculateParentID(id int) int {
	if id <= 10 {
		return 1 // Top-level comments
	}
	if id <= 50 {
		return ((id - 11) % 9) + 2 // Distribute among top-level comments 2-10
	}
	// For deeper levels, use a simple formula
	return ((id - 51) % 40) + 11
}

// formatKidsArray formats a slice of ints as JSON array string
func formatKidsArray(kids []int) string {
	if len(kids) == 0 {
		return "null"
	}
	
	result := "["
	for i, kid := range kids {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d", kid)
	}
	result += "]"
	return result
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BenchmarkOptimizedFetchAllDescendants tests the optimized version with custom rate limiting
func BenchmarkOptimizedFetchAllDescendants(b *testing.B) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	// Override the client with optimized settings
	optimizedClient, err := gohn.NewClientWithOptions(&gohn.ClientOptions{
		HTTPClient: client.GetHTTPClient(),
		RateLimit:  10 * time.Millisecond, // Much faster for testing
		BurstSize:  100,
	})
	if err != nil {
		b.Fatalf("Failed to create optimized client: %v", err)
	}
	
	// Use the same BaseURL as the test setup
	optimizedClient.BaseURL = client.BaseURL

	// Setup a simple story with 10 comments
	mockStory := &gohn.Item{
		ID:          intPtr(1),
		Type:        stringPtr("story"),
		Kids:        &[]int{2, 3, 4},
		Descendants: intPtr(10),
	}

	// Mock responses for baseline test
	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 10}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7, 8], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "kids": [9, 10, 11], "parent": 1}`)
	})
	for i := 5; i <= 11; i++ {
		parent := 2
		if i > 6 {
			parent = 3
		}
		if i > 8 {
			parent = 4
		}
		localI := i 
		localParent := parent 
		mux.HandleFunc(fmt.Sprintf("/item/%d.json", localI), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"id": %d, "type": "comment", "parent": %d}`, localI, localParent)
		})
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		comments, err := optimizedClient.Items.FetchAllDescendantsOptimized(ctx, mockStory, nil, 5)
		if err != nil {
			b.Fatalf("Failed to fetch comments: %v", err)
		}
		if len(comments) != 10 {
			b.Fatalf("Expected 10 comments, got %d", len(comments))
		}
	}
}

// BenchmarkCurrentPerformanceBaseline runs a simple test to establish baseline performance
// for fetching a smaller number of comments to compare against optimizations
func BenchmarkCurrentPerformanceBaseline(b *testing.B) {
	client, mux, _, teardown := setup.Init()
	defer teardown()

	// Setup a simple story with 10 comments
	mockStory := &gohn.Item{
		ID:          intPtr(1),
		Type:        stringPtr("story"),
		Kids:        &[]int{2, 3, 4},
		Descendants: intPtr(10),
	}

	// Mock responses for baseline test
	mux.HandleFunc("/item/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "type": "story", "kids": [2, 3, 4], "descendants": 10}`)
	})
	mux.HandleFunc("/item/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 2, "type": "comment", "kids": [5, 6], "parent": 1}`)
	})
	mux.HandleFunc("/item/3.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 3, "type": "comment", "kids": [7, 8], "parent": 1}`)
	})
	mux.HandleFunc("/item/4.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 4, "type": "comment", "kids": [9, 10, 11], "parent": 1}`)
	})
	for i := 5; i <= 11; i++ {
		parent := 2
		if i > 6 {
			parent = 3
		}
		if i > 8 {
			parent = 4
		}
		localI := i // Capture for closure
		localParent := parent // Capture for closure
		mux.HandleFunc(fmt.Sprintf("/item/%d.json", localI), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"id": %d, "type": "comment", "parent": %d}`, localI, localParent)
		})
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		comments, err := client.Items.FetchAllDescendants(ctx, mockStory, nil)
		if err != nil {
			b.Fatalf("Failed to fetch comments: %v", err)
		}
		if len(comments) != 10 {
			b.Fatalf("Expected 10 comments, got %d", len(comments))
		}
	}
}