package gohn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	USER_URL = "https://hacker-news.firebaseio.com/v0/user/%s.json"
)

// User represents a single user from the Hacker News API.
type User struct {
	ID        string `json:"id"`
	Created   int    `json:"created"`
	Karma     int    `json:"karma"`
	About     string `json:"about"`
	Submitted []int  `json:"submitted"`
}

// GetUser returns a User given a username.
func GetUser(username string) (User, error) {
	var user User

	url := fmt.Sprintf(USER_URL, username)
	resp, err := http.Get(url)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
