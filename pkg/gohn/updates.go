package gohn

import (
	"context"
)

// UpdatesService handles communication with the updates
// related methods of the Hacker News API.
type UpdatesService service

// Update represents a items and profiles changes from the Hacker News API.
// https://github.com/HackerNews/API#changed-items-and-profiles
type Update struct {
	Items    *[]int    `json:"items"`
	Profiles *[]string `json:"profiles"`
}

// UPDATES_URL is the URL for the updates endpoint.
const (
	UPDATES_URL = "updates.json"
)

// Get items and profiles changes.
// https://github.com/HackerNews/API#changed-items-and-profiles
func (s UpdatesService) Get(ctx context.Context) (*Update, error) {
	req, err := s.client.NewRequest("GET", UPDATES_URL)
	if err != nil {
		return nil, err
	}

	var updates *Update
	_, err = s.client.Do(ctx, req, &updates)
	if err != nil {
		return nil, err
	}

	return updates, nil
}
