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

// FetchAllDescendants returns a map of all the comments for a given Item.
// The map key is the ID of the item. The map value is a pointer to the item itself.
// The map can be used to retrieve the comments for a given item by
// traversing the Kids slice of the item recursively (N-ary tree preorder traversal).
// See an implementation in the example directory.
// If the ItemProcessor returns an error, the item will not be added to the map.
// Its kids will be added to the queue only if the ItemProcessors returns false, together with the error.
// For more information on the ItemProcessor, check the gohn/processors package.
func (s *ItemsService) FetchAllDescendants(ctx context.Context, item *Item, fn ItemProcessor) (ItemsIndex, error) {
	if item == nil {
		return nil, errors.New("item is nil")
	}
	if item.Kids == nil {
		return nil, errors.New("item has no kids")
	}
	var wg sync.WaitGroup
	// channel of items to be added to the map
	var commentsChan chan *Item
	// channel of IDs to be retrieved
	// kids found in the commentsChan will be added to this channel
	// buffered so that initializing the queue doesn't block
	var kidsQueue chan int
	// channel to signaling that the processing is done
	done := make(chan struct{})
	// number of items to fetch and process
	var commentsNumToFetch int
	// map of items to return
	var mapCommentById ItemsIndex

	if item.Descendants != nil && *item.Descendants > 0 {
		commentsNumToFetch = *item.Descendants
		mapCommentById = make(ItemsIndex, commentsNumToFetch)
		commentsChan = make(chan *Item, commentsNumToFetch)
		kidsQueue = make(chan int, commentsNumToFetch)
	} else {
		commentsNumToFetch = len(*item.Kids)
		mapCommentById = make(ItemsIndex)
		commentsChan = make(chan *Item)
		kidsQueue = make(chan int, commentsNumToFetch)
	}

	wg.Add(len(*item.Kids))

	// the following variables are used to signal the start of the processing
	start := make(chan struct{})
	var startOnce sync.Once

	// initialize kidsQueue so that the fetching in the for loop can start
	go func() {
		for _, kid := range *item.Kids {
			kidsQueue <- kid
		}
	}()

	// goroutine to close the done channel when all the items are fetched and processed
	go func() {
		// wait for the first kid to be in the queue
		// in this way, we don't close the done channel before the processing has started
		<-start
		// the first kid has been added to the queue
		// it is safe now to start waiting
		wg.Wait()
		close(done)
	}()
L:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		// fetch the item and send it to commentsChan
		case currentId := <-kidsQueue:
			// once the first kid has been added to the queue, signal the start of the processing
			// this is done to avoid closing the done channel before the processing has started
			startOnce.Do(func() {
				close(start)
			})
			go func(wg *sync.WaitGroup, currentId int) {
				it, err := s.Get(ctx, currentId)
				if err != nil {
					// TODO: add better error handling
					wg.Done()
					return
				}
				if fn != nil {
					excludeKids, err := fn(it, wg)
					if err != nil && excludeKids {
						// TODO: add better error handling
						wg.Done()
						return
					} else if err != nil && !excludeKids {
						if it.Kids != nil {
							wg.Add(len(*it.Kids))
							for _, kid := range *it.Kids {
								kidsQueue <- kid
							}
						}
						wg.Done()
						return
					}
				}
				commentsChan <- it
			}(&wg, currentId)
		// add the item to the map and, if it has any kid,
		// add their IDs to the queue so that they can be fetched
		case comment := <-commentsChan:
			if comment.ID != nil {
				mapCommentById[*comment.ID] = comment
				if comment.Kids != nil && len(*comment.Kids) > 0 {
					wg.Add(len(*comment.Kids))
					go func(comment *Item) {
						for _, kid := range *comment.Kids {
							kidsQueue <- kid
						}
					}(comment)
				}
			}
			wg.Done()
		case <-done:
			break L
		}
	}
	return mapCommentById, nil
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
		return nil, ErrInvalidItem{Message: "item is nil"}
	}
	if *item.Type != "comment" {
		return nil, ErrInvalidItem{Message: "item is not a comment"}
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
		item, _ = s.Get(ctx, *item.Parent)
	}
	return storyId, nil
}
