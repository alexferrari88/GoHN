// Package gohn is a wrapper for the Hacker News API which uses goroutines.
package gohn

import (
	"fmt"
)

const (
	TOP_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	BEST_STORIES_URL = "https://hacker-news.firebaseio.com/v0/beststories.json"
	NEW_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	ASK_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/askstories.json"
	SHOW_STORIES_URL = "https://hacker-news.firebaseio.com/v0/showstories.json"
	JOB_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/jobstories.json"
	ITEM_URL         = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

// GetItem returns an Item given an ID.
func GetItem(id int) (Item, error) {
	return retrieveFromURL[Item](fmt.Sprintf(ITEM_URL, id))
}

// GetTopStories returns up to the top 500 stories on Hacker News.
func GetTopStoriesIDs() ([]int, error) {
	return RetrieveIDs(TOP_STORIES_URL)
}

// GetBestStories returns up to the best 500 stories on Hacker News.
func GetBestStoriesIDs() ([]int, error) {
	return RetrieveIDs(BEST_STORIES_URL)
}

// GetNewStories returns up to the newest 500 stories on Hacker News.
func GetNewStoriesIDs() ([]int, error) {
	return RetrieveIDs(NEW_STORIES_URL)
}

// GetAskStories returns up to 200 of the latest Ask stories on Hacker News.
func GetAskStoriesIDs() ([]int, error) {
	return RetrieveIDs(ASK_STORIES_URL)
}

// GetShowStories returns up to 200 of the latest Show stories on Hacker News.
func GetShowStoriesIDs() ([]int, error) {
	return RetrieveIDs(SHOW_STORIES_URL)
}

// GetJobStories returns up to 200 of the latest Job stories on Hacker News.
func GetJobStoriesIDs() ([]int, error) {
	return RetrieveIDs(JOB_STORIES_URL)
}
