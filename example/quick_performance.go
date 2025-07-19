package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

func main() {
	fmt.Println("Quick Performance Test - Real HN API")
	fmt.Println("=====================================")

	// Create an optimized client for better performance
	// Note: Be respectful with rate limits when using the real API
	optimizedClient, err := gohn.NewClientWithOptions(&gohn.ClientOptions{
		RateLimit: 200 * time.Millisecond, // 5 requests per second - respectful but faster
		BurstSize: 10,                     // Allow small bursts
	})
	if err != nil {
		log.Fatalf("Failed to create optimized client: %v", err)
	}

	ctx := context.Background()

	// Get the HN thread item
	fmt.Println("Fetching story metadata...")
	story, err := optimizedClient.Items.Get(ctx, 34591291)
	if err != nil {
		log.Fatalf("Failed to get HN thread 34591291: %v", err)
	}

	if story.Kids == nil {
		log.Fatalf("Story has no comments to fetch")
	}

	fmt.Printf("Story: %s\n", *story.Title)
	if story.Descendants != nil {
		fmt.Printf("Expected comments: %d\n", *story.Descendants)
	} else {
		fmt.Printf("Direct children: %d\n", len(*story.Kids))
	}
	fmt.Printf("Started at: %v\n", time.Now().Format("15:04:05"))
	fmt.Println()

	// Test the optimized version with limited scope to avoid overwhelming the API
	fmt.Println("Testing optimized FetchAllDescendantsOptimized (first 20 top-level comments only)...")
	
	// Create a limited version of the story for testing
	limitedKids := (*story.Kids)[:min(20, len(*story.Kids))]
	limitedStory := &gohn.Item{
		ID:   story.ID,
		Type: story.Type,
		Kids: &limitedKids,
	}

	start := time.Now()
	comments, err := optimizedClient.Items.FetchAllDescendantsOptimized(ctx, limitedStory, processors.UnescapeHTML(), 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	duration := time.Since(start)
	fmt.Printf("Completed in: %v\n", duration)
	fmt.Printf("Comments fetched: %d\n", len(comments))
	fmt.Printf("Average time per comment: %v\n", duration/time.Duration(len(comments)))
	fmt.Printf("Estimated time for all %d comments: %v\n", 
		*story.Descendants, 
		duration*time.Duration(*story.Descendants)/time.Duration(len(comments)))

	fmt.Println()
	fmt.Println("Performance test completed!")
	fmt.Println("Note: This was a limited test to be respectful of the HN API.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}