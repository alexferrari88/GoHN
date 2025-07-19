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
	fmt.Println("GoHN Performance Test - HN Thread 34591291 with 271 comments")
	fmt.Println("============================================================")

	// Test with optimized configuration for performance
	optimizedClient, err := gohn.NewClientWithOptions(&gohn.ClientOptions{
		RateLimit: 100 * time.Millisecond, // 10 requests per second (be respectful)
		BurstSize: 20,                     // Allow bursts of 20 requests
	})
	if err != nil {
		log.Fatalf("Failed to create optimized client: %v", err)
	}

	// Test with default configuration
	defaultClient, err := gohn.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create default client: %v", err)
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
	fmt.Printf("Expected comments: %d\n", *story.Descendants)
	fmt.Println()

	// Test 1: Default client with original method
	fmt.Println("Test 1: Default client (1 req/sec) with FetchAllDescendants")
	start := time.Now()
	_, err = defaultClient.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("Completed in: %v\n", duration)
	}
	fmt.Println()

	// Test 2: Optimized client with original method
	fmt.Println("Test 2: Optimized client (10 req/sec) with FetchAllDescendants")
	start = time.Now()
	comments1, err := optimizedClient.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("Completed in: %v\n", duration)
		fmt.Printf("Comments fetched: %d\n", len(comments1))
	}
	fmt.Println()

	// Test 3: Optimized client with optimized method
	fmt.Println("Test 3: Optimized client (10 req/sec) with FetchAllDescendantsOptimized")
	start = time.Now()
	comments2, err := optimizedClient.Items.FetchAllDescendantsOptimized(ctx, story, processors.UnescapeHTML(), 0) // Auto-determine worker count
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("Completed in: %v\n", duration)
		fmt.Printf("Comments fetched: %d\n", len(comments2))
	}
	fmt.Println()

	// Test 4: Optimized client with more aggressive settings (be careful with real API!)
	fmt.Println("Test 4: Very optimized client (20 req/sec) - USE WITH CAUTION!")
	veryOptimizedClient, err := gohn.NewClientWithOptions(&gohn.ClientOptions{
		RateLimit: 50 * time.Millisecond, // 20 requests per second
		BurstSize: 50,                    // Allow larger bursts
	})
	if err != nil {
		log.Fatalf("Failed to create very optimized client: %v", err)
	}

	start = time.Now()
	comments3, err := veryOptimizedClient.Items.FetchAllDescendantsOptimized(ctx, story, processors.UnescapeHTML(), 25) // Use 25 workers
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("Completed in: %v\n", duration)
		fmt.Printf("Comments fetched: %d\n", len(comments3))
	}

	fmt.Println()
	fmt.Println("Performance test completed!")
	fmt.Println("Note: Be respectful of the HN API when using higher rate limits.")
}