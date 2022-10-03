package gohn

import "encoding/json"

// Update represents a items and profiles changes from the Hacker News API.
// https://github.com/HackerNews/API#changed-items-and-profiles
type Update struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

// UPDATES_URL is the URL for the updates endpoint.
const (
	UPDATES_URL = "https://hacker-news.firebaseio.com/v0/updates.json"
)

// GetUpdates returns a slice of IDs for the given URL.
func (c client) GetUpdates() (Update, error) {
	var updates Update

	resp, err := c.retrieveFromURL(UPDATES_URL)
	if err != nil {
		return updates, err
	}
	err = json.Unmarshal(resp, &updates)
	if err != nil {
		return updates, err
	}

	return updates, nil
}
