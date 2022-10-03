package gohn

import "encoding/json"

// Update represents a items and profiles changes from the Hacker News API.
// https://github.com/HackerNews/API#changed-items-and-profiles
type Update struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

// updates_url is the URL for the updates endpoint.
const (
	updates_url = "https://hacker-news.firebaseio.com/v0/updates.json"
)

// GetUpdates returns a slice of IDs for the given URL.
func (c client) GetUpdates() (Update, error) {
	var updates Update

	resp, err := c.retrieveFromURL(updates_url)
	if err != nil {
		return updates, err
	}
	err = json.Unmarshal(resp, &updates)
	if err != nil {
		return updates, err
	}

	return updates, nil
}
