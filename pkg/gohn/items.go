package gohn

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ITEM_URL is the URL to retrieve an item given an ID.
// MAX_ITEM_ID_URL is the URL that retrieves the current largest item id.
const (
	ITEM_URL        = "item/%d.json"
	MAX_ITEM_ID_URL = "maxitem.json"
)

// ItemService handles retrieving items from the Hacker News API.
type ItemsService service

// Item represents a single item from the Hacker News API.
// https://github.com/HackerNews/API#items
type Item struct {
	ID          *int    `json:"id,omitempty"`
	Deleted     *bool   `json:"deleted,omitempty"`
	Type        *string `json:"type,omitempty"`
	By          *string `json:"by,omitempty"`
	Time        *int    `json:"time,omitempty"`
	Text        *string `json:"text,omitempty"`
	Dead        *bool   `json:"dead,omitempty"`
	Parent      *int    `json:"parent,omitempty"`
	Poll        *int    `json:"poll,omitempty"`
	Kids        *[]int  `json:"kids,omitempty"`
	Position    *int
	URL         *string `json:"url,omitempty"`
	Score       *int    `json:"score,omitempty"`
	Title       *string `json:"title,omitempty"`
	Parts       *[]int  `json:"parts,omitempty"`
	Descendants *int    `json:"descendants,omitempty"`
}

// ItemProcessor is used by ItemsService.Get and ItemsService.FetchAllKids
// to process items after they are retrieved.
// The package itemprocessor provides some common implementations.
type ItemProcessor func(*Item, *sync.WaitGroup) (bool, error)

// Get returns an Item given an ID.
func (s *ItemsService) Get(ctx context.Context, id int) (*Item, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf(ITEM_URL, id))
	if err != nil {
		return nil, err
	}

	var item *Item
	_, err = s.client.Do(ctx, req, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetIDsFromURL returns a slice of Items' IDs for the given URL.
func (s *ItemsService) GetIDsFromURL(ctx context.Context, url string) ([]*int, error) {
	req, err := s.client.NewRequest("GET", url)
	if err != nil {
		return nil, err
	}
	var ids []*int
	_, err = s.client.Do(ctx, req, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// commentFetchConfig holds the configuration and channels for comment fetching.
// This struct encapsulates all the state needed for concurrent comment retrieval,
// providing a clean separation of concerns and making the code more testable.
type commentFetchConfig struct {
	commentsChan   chan *Item    // Channel for passing fetched comments to the result processor
	kidsQueue      chan int      // Buffered channel for queuing comment IDs to be fetched
	results        ItemsIndex    // Thread-safe map for storing fetched comments by ID
	expectedCount  int           // Expected number of comments (used for buffer sizing optimization)
}

// initializeCommentFetching sets up channels and buffers for the comment fetching process.
// This function optimizes channel buffer sizes based on the expected number of descendants:
// - If item.Descendants is available, it pre-allocates buffers for better performance
// - If not available, it falls back to using the number of direct children
// This optimization reduces memory allocations and improves concurrent performance.
func (s *ItemsService) initializeCommentFetching(item *Item) *commentFetchConfig {
	var expectedCount int
	var commentsChan chan *Item
	var kidsQueue chan int
	var results ItemsIndex

	if item.Descendants != nil && *item.Descendants > 0 {
		expectedCount = *item.Descendants
		results = make(ItemsIndex, expectedCount)
		// Pre-allocate larger buffers for better performance
		bufferSize := min(expectedCount, 1000) // Cap buffer at 1000 to avoid excessive memory usage
		commentsChan = make(chan *Item, bufferSize)
		kidsQueue = make(chan int, bufferSize)
	} else {
		expectedCount = len(*item.Kids)
		results = make(ItemsIndex, expectedCount*2) // Pre-allocate more space for nested comments
		bufferSize := min(expectedCount*10, 500) // Estimate more children than direct kids
		commentsChan = make(chan *Item, bufferSize)
		kidsQueue = make(chan int, bufferSize)
	}

	return &commentFetchConfig{
		commentsChan:  commentsChan,
		kidsQueue:     kidsQueue,
		results:       results,
		expectedCount: expectedCount,
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


// FetchAllDescendants returns a map of all the comments for a given Item.
// The map key is the ID of the item. The map value is a pointer to the item itself.
// The map can be used to retrieve the comments for a given item by
// traversing the Kids slice of the item recursively (N-ary tree preorder traversal).
// See an implementation in the example directory.
//
// CONCURRENT ARCHITECTURE:
// This function implements a sophisticated concurrent comment fetching system:
// 1. Uses a work queue pattern with channels for distributing fetch tasks
// 2. Employs a WaitGroup to track in-progress work and signal completion
// 3. Processes results dynamically, adding new work items as comments are fetched
// 4. Provides graceful error handling and context cancellation support
//
// PROCESSOR INTEGRATION:
// If the ItemProcessor returns an error, the item will not be added to the map.
// Its kids will be added to the queue only if the ItemProcessor returns false, together with the error.
// For more information on the ItemProcessor, check the gohn/processors package.
//
// PERFORMANCE OPTIMIZATIONS:
// - Pre-allocates channel buffers based on expected comment count
// - Uses separate goroutines for different aspects of the work
// - Minimizes lock contention through channel-based communication
func (s *ItemsService) FetchAllDescendants(ctx context.Context, item *Item, fn ItemProcessor) (ItemsIndex, error) {
	if item == nil {
		return nil, errors.New("item is nil")
	}
	if item.Kids == nil {
		return nil, errors.New("item has no kids")
	}

	// Initialize configuration for the fetch operation
	config := s.initializeCommentFetching(item)
	
	// Use a WaitGroup to track work in progress
	var wg sync.WaitGroup
	
	// Error channel for goroutine error reporting
	errChan := make(chan error, 1)
	
	// Channel to signal completion
	done := make(chan struct{})
	
	// Add initial work items
	wg.Add(len(*item.Kids))
	
	// Start initial comment fetchers
	go s.populateKidsQueue(config.kidsQueue, *item.Kids)
	
	// Start the completion monitor
	go s.monitorCompletion(&wg, done)
	
	// Process work items until completion
	return s.processWorkQueue(ctx, config, &wg, errChan, done, fn)
}

// FetchAllDescendantsOptimized is an optimized version of FetchAllDescendants
// that uses configurable worker pools and batching for better performance.
// This function is recommended for large comment trees (>50 comments).
func (s *ItemsService) FetchAllDescendantsOptimized(ctx context.Context, item *Item, fn ItemProcessor, workerCount int) (ItemsIndex, error) {
	if item == nil {
		return nil, errors.New("item is nil")
	}
	if item.Kids == nil {
		return nil, errors.New("item has no kids")
	}

	if workerCount <= 0 {
		// Determine optimal worker count based on expected descendants
		if item.Descendants != nil && *item.Descendants > 0 {
			// For large comment trees, use more workers but cap it
			workerCount = min(max(*item.Descendants/20, 5), 50)
		} else {
			// Default worker count for smaller trees
			workerCount = min(len(*item.Kids)*2, 10)
		}
	}

	// Initialize configuration for the fetch operation
	config := s.initializeCommentFetching(item)
	
	// Use a WaitGroup to track work in progress
	var wg sync.WaitGroup
	
	// Error channel for goroutine error reporting
	errChan := make(chan error, 1)
	
	// Channel to signal completion
	done := make(chan struct{})
	
	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go s.worker(ctx, config, &wg, errChan, fn)
	}
	
	// Add initial work items
	wg.Add(len(*item.Kids))
	
	// Start initial comment fetchers
	go s.populateKidsQueue(config.kidsQueue, *item.Kids)
	
	// Start the completion monitor
	go s.monitorCompletion(&wg, done)
	
	// Process results until completion
	return s.processWorkQueueOptimized(ctx, config, &wg, errChan, done)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// populateKidsQueue populates the kids queue with initial IDs.
// This function runs in a separate goroutine to avoid blocking the main processing loop.
// It's used both for initial population and for adding newly discovered child comments.
func (s *ItemsService) populateKidsQueue(kidsQueue chan int, kids []int) {
	for _, kid := range kids {
		kidsQueue <- kid
	}
}

// monitorCompletion waits for all work to complete and signals done.
// This function runs in a separate goroutine to monitor the WaitGroup and
// signal completion when all concurrent fetch operations have finished.
// This approach separates completion monitoring from the main processing loop.
func (s *ItemsService) monitorCompletion(wg *sync.WaitGroup, done chan struct{}) {
	wg.Wait()
	close(done)
}

// processWorkQueue handles the main processing loop with channels.
// This is the heart of the concurrent processing system, implementing a select-based
// event loop that coordinates between:
// - Context cancellation (for graceful shutdown)
// - Error propagation from worker goroutines
// - Work item distribution to fetch goroutines
// - Result collection and processing
// - Completion signaling when all work is done
func (s *ItemsService) processWorkQueue(ctx context.Context, config *commentFetchConfig, wg *sync.WaitGroup, errChan chan error, done chan struct{}, fn ItemProcessor) (ItemsIndex, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errChan:
			return nil, err
		case itemID := <-config.kidsQueue:
			go s.fetchAndProcessComment(ctx, itemID, config, wg, errChan, fn)
		case comment := <-config.commentsChan:
			s.handleCommentResult(comment, config, wg)
		case <-done:
			return config.results, nil
		}
	}
}

// fetchAndProcessComment fetches a single comment and processes it.
// This function runs in a separate goroutine for each comment being fetched, providing
// concurrent processing. It handles:
// - HTTP request to fetch the comment data
// - ItemProcessor integration with proper error handling
// - Result forwarding to the main processing loop
// - WaitGroup coordination for completion tracking
// - Error propagation through the error channel
func (s *ItemsService) fetchAndProcessComment(ctx context.Context, itemID int, config *commentFetchConfig, wg *sync.WaitGroup, errChan chan error, fn ItemProcessor) {
	defer wg.Done()
	
	item, err := s.Get(ctx, itemID)
	if err != nil {
		select {
		case errChan <- err:
		default:
		}
		return
	}

	// Process with the provided processor if available
	if fn != nil {
		excludeKids, err := fn(item, wg)
		if err != nil {
			if excludeKids {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			// If processor returns error but doesn't want to exclude kids,
			// we still add the kids to the queue but don't add this item to results
			if item.Kids != nil && len(*item.Kids) > 0 {
				wg.Add(len(*item.Kids))
				go s.populateKidsQueue(config.kidsQueue, *item.Kids)
			}
			return
		}
	}

	// Send the item to be processed
	config.commentsChan <- item
}

// handleCommentResult processes a comment result and adds its kids to the queue.
// This function is called from the main processing loop when a comment is successfully
// fetched and processed. It:
// - Stores the comment in the results map
// - Discovers any child comments and adds them to the work queue
// - Updates the WaitGroup to track the newly discovered work
// This enables the dynamic discovery and fetching of the complete comment tree.
func (s *ItemsService) handleCommentResult(comment *Item, config *commentFetchConfig, wg *sync.WaitGroup) {
	if comment.ID != nil {
		config.results[*comment.ID] = comment
		if comment.Kids != nil && len(*comment.Kids) > 0 {
			wg.Add(len(*comment.Kids))
			go s.populateKidsQueue(config.kidsQueue, *comment.Kids)
		}
	}
}

// worker is a worker goroutine that processes items from the queue.
// This is used by the optimized version to maintain a fixed pool of workers
// instead of creating a new goroutine for each item.
func (s *ItemsService) worker(ctx context.Context, config *commentFetchConfig, wg *sync.WaitGroup, errChan chan error, fn ItemProcessor) {
	for {
		select {
		case <-ctx.Done():
			return
		case itemID, ok := <-config.kidsQueue:
			if !ok {
				return // Channel closed
			}
			s.fetchAndProcessComment(ctx, itemID, config, wg, errChan, fn)
		}
	}
}

// processWorkQueueOptimized handles the optimized processing loop with worker pools.
// This version separates the work distribution from result collection for better performance.
func (s *ItemsService) processWorkQueueOptimized(ctx context.Context, config *commentFetchConfig, wg *sync.WaitGroup, errChan chan error, done chan struct{}) (ItemsIndex, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errChan:
			return nil, err
		case comment := <-config.commentsChan:
			s.handleCommentResult(comment, config, wg)
		case <-done:
			close(config.kidsQueue) // Signal workers to stop
			return config.results, nil
		}
	}
}

// GetMaxID returns the ID of the most recent item.
// https://github.com/HackerNews/API#max-item-id
func (s *ItemsService) GetMaxID(ctx context.Context) (*int, error) {
	req, err := s.client.NewRequest("GET", MAX_ITEM_ID_URL)
	if err != nil {
		return nil, err
	}
	var maxID *int
	_, err = s.client.Do(ctx, req, &maxID)
	if err != nil {
		return nil, err
	}
	return maxID, nil
}

// GetStoryIdFromComment returns the ID of the story for a given comment.
func (s *ItemsService) GetStoryIdFromComment(ctx context.Context, item *Item) (*int, error) {
	if item == nil {
		return nil, InvalidItemError{Message: "item is nil"}
	}
	if *item.Type != "comment" {
		return nil, InvalidItemError{Message: "item is not a comment"}
	}
	var storyId *int
	for {
		if *item.Type == "story" {
			storyId = item.ID
			break
		}
		if item.Parent == nil {
			break
		}
		parentItem, err := s.Get(ctx, *item.Parent)
		if err != nil {
			return nil, err
		}
		item = parentItem
	}
	return storyId, nil
}