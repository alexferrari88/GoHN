package gohn

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

type ItemWithComments struct {
	Item
	Comments []ItemWithComments `json:"comments"`
}

// RetrieveIDs returns a slice of IDs for the given URL.
func RetrieveIDs(url string) ([]int, error) {
	var ids []int

	resp, err := http.Get(url)
	if err != nil {
		return ids, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ids, err
	}

	err = json.Unmarshal(body, &ids)
	if err != nil {
		return ids, err
	}

	return ids, nil
}

func GetItemWithComments(id int) (ItemWithComments, error) {
	var item ItemWithComments

	item.Item, _ = GetItem(id)

	for _, commentID := range item.Kids {
		comment, _ := GetItemWithComments(commentID)
		item.Comments = append(item.Comments, comment)
	}

	return item, nil
}
