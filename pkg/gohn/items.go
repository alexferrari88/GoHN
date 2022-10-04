package gohn

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// MAX_ITEM_ID_URL is the URL that retrieves the current largest item id.
const (
	ITEM_URL        = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	MAX_ITEM_ID_URL = "https://hacker-news.firebaseio.com/v0/maxitem.json"
)

// Item represents a single item from the Hacker News API.
// https://github.com/HackerNews/API#items
type Item struct {
	ID          int    `json:"id"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int    `json:"time"`
	Text        string `json:"text"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Poll        int    `json:"poll"`
	Kids        []int  `json:"kids"`
	Order       int
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

// ItemsIndex is a map of items indexed by their ID.
type ItemsIndex map[int]Item

// ItemProcessor is used by GetItem and Item.RetrieveKidsItems
// to process items after they are retrieved.
// The package itemprocessor provides some common implementations.
type ItemProcessor func(*Item) error

// GetItem returns an Item given an ID.
func (c client) GetItem(id int) (Item, error) {
	var item Item
	resp, err := c.retrieveFromURL(fmt.Sprintf(ITEM_URL, id))
	if err != nil {
		return item, err
	}
	err = json.Unmarshal(resp, &item)
	if err != nil {
		return item, err
	}
	return item, nil
}

// RetrieveIDs returns a slice of IDs for the given URL.
func (c client) RetrieveIDs(url string) ([]int, error) {
	var ids []int
	resp, err := c.retrieveFromURL(url)
	if err != nil {
		return ids, err
	}
	err = json.Unmarshal(resp, &ids)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

// RetrieveKidsItems returns a map of all the comments for a given Item.
// The map key is the ID of the item. The map value is the item itself.
// The map can be used to retrieve the comments for a given item by
// traversing the Kids slice of the item recursively (N-ary tree preorder traversal).
// See an implementation in the example directory.
// If the ItemProcessor returns an error, the item will not be added to the map.
func (c client) RetrieveKidsItems(item Item, fn ItemProcessor) ItemsIndex {
	// the map of items to return
	mapCommentById := make(ItemsIndex)
	// the channel of items to be added to the map
	commentsChan := make(chan Item)
	// the channel of IDs to be retrieved
	// kids found in the commentsChan will be added to this channel
	// buffered so that initializing the queue doesn't block
	kidsQueue := make(chan int, len(item.Kids))
	// Use an atomic counter to keep track of the number of items
	commentsNumToFetch := int32(len(item.Kids))
	// initialize kidsQueue so that the fetching in the for loop can start
	for _, kid := range item.Kids {
		kidsQueue <- kid
	}
L:
	for {
		select {
		case currentId := <-kidsQueue:
			if commentsNumToFetch > 0 {
				go func() {
					it, err := c.GetItem(currentId)
					if err != nil {
						// TODO: add better error handling
						atomic.AddInt32(&commentsNumToFetch, -1)
						return
					}
					if fn != nil {
						err = fn(&it)
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
			if comment.ID != 0 {
				mapCommentById[comment.ID] = comment
				atomic.AddInt32(&commentsNumToFetch, int32(len(comment.Kids)))
				go func() {
					for _, kid := range comment.Kids {
						kidsQueue <- kid
					}
				}()
			}
		default:
			if commentsNumToFetch == 0 {
				break L
			}
		}
	}
	return mapCommentById
}

// GetMaxItemID returns the ID of the most recent item.
func (c client) GetMaxItemID() (int, error) {
	var id int
	resp, err := c.retrieveFromURL(MAX_ITEM_ID_URL)
	if err != nil {
		return id, err
	}
	err = json.Unmarshal(resp, &id)
	if err != nil {
		return id, err
	}
	return id, nil
}

// IsTopLevelComment checks if an Item is a top level comment in a story.
func (s Story) IsTopLevelComment(item Item) bool {
	for _, kid := range s.StoryItem.Kids {
		if kid == item.ID {
			return true
		}
	}
	return false
}
