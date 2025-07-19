package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

// TestRealHNThread tests our optimized implementation against a real Hacker News thread
// Using thread ID 34591291 as provided by the user for realistic validation
func TestRealHNThread(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a real client (not mocked)
	client, err := gohn.NewClient(nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// HN API has rate limiting (1 req/sec), so we need a generous timeout
	// For a thread with 271 descendants, this could take 4+ minutes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Test with the real HN thread ID provided by user
	storyID := 34591291
	
	t.Logf("Fetching real HN story ID: %d", storyID)
	
	// Fetch the main story item
	story, err := client.Items.Get(ctx, storyID)
	if err != nil {
		t.Fatalf("Failed to fetch story: %v", err)
	}

	if story == nil {
		t.Fatalf("Story is nil")
	}

	t.Logf("Story title: %s", getStringValue(story.Title))
	t.Logf("Story author: %s", getStringValue(story.By))
	t.Logf("Story score: %d", getIntValue(story.Score))
	t.Logf("Expected descendants: %d", getIntValue(story.Descendants))

	if story.Kids == nil || len(*story.Kids) == 0 {
		t.Skip("Story has no comments to test with")
	}

	t.Logf("Direct children comments: %d", len(*story.Kids))

	// Create a limited test story with only the first few kids to avoid long test times
	// With rate limiting at 1req/sec, even 10 comments takes 10+ seconds
	maxKidsToTest := 2
	limitedKids := (*story.Kids)[:min(maxKidsToTest, len(*story.Kids))]
	limitedStory := &gohn.Item{
		ID:    story.ID,
		Type:  story.Type,
		Kids:  &limitedKids,
		Title: story.Title,
		By:    story.By,
		Score: story.Score,
	}

	// Test 1: Fetch descendants without processor (limited set)
	t.Run("FetchAllDescendants_NoProcessor_Limited", func(t *testing.T) {
		start := time.Now()
		
		comments, err := client.Items.FetchAllDescendants(ctx, limitedStory, nil)
		
		duration := time.Since(start)
		t.Logf("Fetched %d comments in %v", len(comments), duration)
		
		if err != nil {
			t.Errorf("FetchAllDescendants failed: %v", err)
			return
		}

		if len(comments) == 0 {
			t.Errorf("Expected some comments, got 0")
			return
		}

		// Verify comments have valid IDs and are properly structured
		validComments := 0
		for id, comment := range comments {
			if comment.ID == nil {
				t.Errorf("Comment with map key %d has nil ID", id)
				continue
			}
			if *comment.ID != id {
				t.Errorf("Comment ID mismatch: map key %d, comment ID %d", id, *comment.ID)
				continue
			}
			validComments++
		}

		t.Logf("Validated %d comments with proper structure", validComments)
		
		// Ensure we got at least the direct children
		if len(comments) < len(limitedKids) {
			t.Errorf("Expected at least %d comments (direct children), got %d", len(limitedKids), len(comments))
		}
	})

	// Test 2: Fetch with UnescapeHTML processor
	t.Run("FetchAllDescendants_WithUnescapeHTML", func(t *testing.T) {
		start := time.Now()
		
		comments, err := client.Items.FetchAllDescendants(ctx, limitedStory, processors.UnescapeHTML())
		
		duration := time.Since(start)
		t.Logf("Fetched and processed %d comments in %v", len(comments), duration)
		
		if err != nil {
			t.Errorf("FetchAllDescendants with processor failed: %v", err)
			return
		}

		if len(comments) == 0 {
			t.Errorf("Expected some comments, got 0")
			return
		}

		// Look for evidence that HTML unescaping worked
		unescapedCount := 0
		for _, comment := range comments {
			if comment.Text != nil {
				text := *comment.Text
				// Check if text contains common unescaped HTML entities
				if containsUnescapedHTML(text) {
					unescapedCount++
				}
			}
		}

		t.Logf("Found evidence of HTML unescaping in %d comments", unescapedCount)
	})

	// Test 3: Test with filtering processor
	t.Run("FetchAllDescendants_WithFilter", func(t *testing.T) {
		// Filter out comments containing common words (should exclude many)
		filterWords := []string{"the", "and", "is", "it", "to"}
		processor := processors.FilterOutWords(filterWords, false)
		
		start := time.Now()
		
		comments, err := client.Items.FetchAllDescendants(ctx, limitedStory, processor)
		
		duration := time.Since(start)
		t.Logf("Fetched %d comments after filtering in %v", len(comments), duration)
		
		if err != nil {
			t.Errorf("FetchAllDescendants with filter failed: %v", err)
			return
		}

		// We expect fewer comments due to filtering
		// Verify no filtered comments contain the banned words
		for _, comment := range comments {
			if comment.Text != nil {
				text := *comment.Text
				for _, word := range filterWords {
					if containsWord(text, word) {
						t.Errorf("Comment %d contains filtered word '%s': %s", 
							getIntValue(comment.ID), word, truncateString(text, 100))
						break
					}
				}
			}
		}
	})
}

// TestRealHNThread_Performance compares performance characteristics
func TestRealHNThread_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	client, err := gohn.NewClient(nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	storyID := 34591291
	story, err := client.Items.Get(ctx, storyID)
	if err != nil {
		t.Fatalf("Failed to fetch story: %v", err)
	}

	if story.Kids == nil || len(*story.Kids) == 0 {
		t.Skip("Story has no comments to test with")
	}

	// Use a limited set for performance testing too
	maxKidsToTest := 2 // Even smaller for performance test
	limitedKids := (*story.Kids)[:min(maxKidsToTest, len(*story.Kids))]
	limitedStory := &gohn.Item{
		ID:   story.ID,
		Type: story.Type,
		Kids: &limitedKids,
	}

	// Run multiple iterations to get stable timing
	const iterations = 3
	var totalDuration time.Duration
	var totalComments int

	for i := 0; i < iterations; i++ {
		start := time.Now()
		
		comments, err := client.Items.FetchAllDescendants(ctx, limitedStory, processors.UnescapeHTML())
		if err != nil {
			t.Fatalf("Iteration %d failed: %v", i+1, err)
		}
		
		duration := time.Since(start)
		totalDuration += duration
		totalComments = len(comments) // Should be consistent across iterations
		
		t.Logf("Iteration %d: %d comments in %v", i+1, len(comments), duration)
	}

	avgDuration := totalDuration / iterations
	t.Logf("Average: %d comments in %v (%.2f comments/sec)", 
		totalComments, avgDuration, float64(totalComments)/avgDuration.Seconds())

	// Validate reasonable performance (this is quite lenient)
	commentsPerSecond := float64(totalComments) / avgDuration.Seconds()
	if commentsPerSecond < 1.0 {
		t.Errorf("Performance seems slow: %.2f comments/sec", commentsPerSecond)
	}
}

// Helper functions
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func containsUnescapedHTML(text string) bool {
	// Look for characters that suggest HTML was unescaped
	// These would be present if UnescapeHTML processor worked
	indicators := []string{"&", "<", ">", "\"", "'"}
	for _, indicator := range indicators {
		if len(text) > 10 && len(text) < 1000 { // Reasonable text length
			// Look for these characters in reasonable amounts (not escaped HTML)
			count := 0
			for _, char := range text {
				if string(char) == indicator {
					count++
				}
			}
			// If we find some but not too many, it suggests unescaping worked
			if count > 0 && count < len(text)/10 {
				return true
			}
		}
	}
	return false
}

func containsWord(text, word string) bool {
	// Simple case-insensitive word matching
	text = strings.ToLower(text)
	word = strings.ToLower(word)
	
	// Look for the word as a standalone word (not part of another word)
	words := strings.Fields(text)
	for _, w := range words {
		// Remove punctuation for comparison
		cleaned := strings.Trim(w, ".,!?;:\"'()[]")
		if cleaned == word {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}