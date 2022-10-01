package gohn

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// MAX_ITEM_ID_URL is the URL that retrieves the current largest item id.
const (
	MAX_ITEM_ID_URL = "https://hacker-news.firebaseio.com/v0/maxitem.json"
)

// Item represents a single item from the Hacker News API.
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

// ItemWithComments represents an item with all its comments.
type ItemWithComments struct {
	Item
	Comments []ItemWithComments `json:"comments"`
}

// RetrieveIDs returns a slice of IDs for the given URL.
func RetrieveIDs(url string) ([]int, error) {
	return retrieveFromURL[[]int](url)
}

// GetItemWithComments returns an item with all its comments.
func GetItemWithComments(id int) (ItemWithComments, error) {
	var item ItemWithComments

	item.Item, _ = GetItem(id)

	for _, commentID := range item.Kids {
		comment, _ := GetItemWithComments(commentID)
		item.Comments = append(item.Comments, comment)
	}

	return item, nil
}

// GetMaxItemID returns the ID of the most recent item.
func GetMaxItemID() (int, error) {
	return retrieveFromURL[int](MAX_ITEM_ID_URL)
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
