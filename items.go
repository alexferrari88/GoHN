package gohn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// max_item_id_url is the URL that retrieves the current largest item id.
const (
	max_item_id_url = "https://hacker-news.firebaseio.com/v0/maxitem.json"
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
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

// ItemProcessor is used by GetItem and Item.RetrieveKidsItems to process items after they are retrieved.
// the package itemprocessor provides some common implementations.
type ItemProcessor func(*Item) error

// GetItem returns an Item given an ID.
func GetItem(id int) (Item, error) {
	item, err := retrieveFromURL[Item](fmt.Sprintf(item_url, id))
	if err != nil {
		return item, err
	}
	return item, nil
}

// RetrieveIDs returns a slice of IDs for the given URL.
func RetrieveIDs(url string) ([]int, error) {
	return retrieveFromURL[[]int](url)
}

// RetrieveKidsItems returns all the comments for a given item.
// If the ItemProcessor returns an error, the item will not be added to the map.
func (i *Item) RetrieveKidsItems(fn ItemProcessor) map[int]Item {
	mapCommentById := make(map[int]Item)
	commentsChan := make(chan Item)
	// buffered so that initializing the queue doesn't block
	kidsQueue := make(chan int, len(i.Kids))
	commentsNumToFetch := len(i.Kids)
	// initialize kidsQueue so that the fetching in the for loop can start
	for _, kid := range i.Kids {
		kidsQueue <- kid
	}
L:
	for {
		select {
		case currentId := <-kidsQueue:
			if commentsNumToFetch > 0 {
				go func() {
					it, err := GetItem(currentId)
					if err != nil {
						// TODO: add better error handling
						commentsNumToFetch--
						return
					}
					if fn != nil {
						err = fn(&it)
						if err != nil {
							// TODO: add better error handling
							commentsNumToFetch--
							return
						}
					}
					commentsChan <- it
				}()
			} else {
				break L
			}
		case comment := <-commentsChan:
			commentsNumToFetch--
			if comment.ID != 0 {
				mapCommentById[comment.ID] = comment
				commentsNumToFetch += len(comment.Kids)
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
func GetMaxItemID() (int, error) {
	return retrieveFromURL[int](max_item_id_url)
}

// retrieveFromURL sends a GET request to the given URL and returns unmarsheled values and an error.
func retrieveFromURL[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
