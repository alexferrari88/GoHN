package gohn

import (
	"encoding/json"
	"fmt"
)

// USER_URL is the URL for the user endpoint.
const (
	USER_URL = "https://hacker-news.firebaseio.com/v0/user/%s.json"
)

// User represents a single user from the Hacker News API.
// https://github.com/HackerNews/API#users
type User struct {
	ID        string `json:"id"`
	Created   int    `json:"created"`
	Karma     int    `json:"karma"`
	About     string `json:"about"`
	Submitted []int  `json:"submitted"`
}

// GetUser returns a User given a username.
func (c client) GetUser(username string) (User, error) {
	var user User

	url := fmt.Sprintf(USER_URL, username)
	resp, err := c.retrieveFromURL(url)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
