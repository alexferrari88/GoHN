package gohn

import (
	"context"
	"fmt"
)

// USER_URL is the URL for the user endpoint.
const (
	USER_URL = "user/%s.json"
)

// UsersService handles communication with the user
// related methods of the Hacker News API.
type UsersService service

// User represents a single user from the Hacker News API.
// https://github.com/HackerNews/API#users
type User struct {
	ID        *string `json:"id"`
	Created   *int    `json:"created"`
	Karma     *int    `json:"karma"`
	About     *string `json:"about"`
	Submitted *[]int  `json:"submitted"`
}

// GetByUsername returns a User given a username.
func (s *UsersService) GetByUsername(ctx context.Context, username string) (*User, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf(USER_URL, username))
	if err != nil {
		return nil, err
	}

	var user *User
	_, err = s.client.Do(ctx, req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
