package gohn

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
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

// ItemsIndex is a map of Items indexed by their ID.
type ItemsIndex map[int]*Item

// ItemProcessor is used by ItemsService.Get and ItemsService.FetchAllKids
// to process items after they are retrieved.
// The package itemprocessor provides some common implementations.
type ItemProcessor func(*Item) error

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

// FetchAllKids returns a map of all the comments for a given Item.
// The map key is the ID of the item. The map value is a pointer to the item itself.
// The map can be used to retrieve the comments for a given item by
// traversing the Kids slice of the item recursively (N-ary tree preorder traversal).
// See an implementation in the example directory.
// If the ItemProcessor returns an error, the item will not be added to the map.
func (s *ItemsService) FetchAllKids(ctx context.Context, item *Item, fn ItemProcessor) (ItemsIndex, error) {
	if item == nil {
		return nil, errors.New("item is nil")
	}
	if item.Kids == nil {
		return nil, errors.New("item has no kids")
	}
	// the map of items to return
	mapCommentById := make(ItemsIndex)
	// the channel of items to be added to the map
	commentsChan := make(chan *Item)
	// the channel of IDs to be retrieved
	// kids found in the commentsChan will be added to this channel
	// buffered so that initializing the queue doesn't block
	kidsQueue := make(chan int, len(*item.Kids))
	// Use an atomic counter to keep track of the number of items
	commentsNumToFetch := int32(len(*item.Kids))
	// initialize kidsQueue so that the fetching in the for loop can start
	for _, kid := range *item.Kids {
		kidsQueue <- kid
	}
L:
	for {
		select {
		case currentId := <-kidsQueue:
			if commentsNumToFetch > 0 {
				go func() {
					it, err := s.Get(ctx, currentId)
					if err != nil {
						// TODO: add better error handling
						atomic.AddInt32(&commentsNumToFetch, -1)
						return
					}
					if fn != nil {
						err = fn(it)
						if err != nil {
							// TODO: add better error handling
							atomic.AddInt32(&commentsNumToFetch, -1)
							return
						}
					}
					commentsChan <- it
				}()
			} else {
				break L
			}
		case comment := <-commentsChan:
			atomic.AddInt32(&commentsNumToFetch, -1)
			if comment.ID != nil {
				mapCommentById[*comment.ID] = comment
				if comment.Kids != nil {
					atomic.AddInt32(&commentsNumToFetch, int32(len(*comment.Kids)))
					go func() {
						for _, kid := range *comment.Kids {
							kidsQueue <- kid
						}
					}()
				}
			}
		default:
			if commentsNumToFetch == 0 {
				break L
			}
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
