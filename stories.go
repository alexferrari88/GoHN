// Package gohn is a wrapper for the Hacker News API which uses goroutines.
package gohn

import (
	"fmt"
)

const (
	top_stories_url  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	best_stories_url = "https://hacker-news.firebaseio.com/v0/beststories.json"
	new_stories_url  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	ask_stories_url  = "https://hacker-news.firebaseio.com/v0/askstories.json"
	show_stories_url = "https://hacker-news.firebaseio.com/v0/showstories.json"
	job_stories_url  = "https://hacker-news.firebaseio.com/v0/jobstories.json"
	item_url         = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

// GetItem returns an Item given an ID.
func GetItem(id int) (Item, error) {
	return retrieveFromURL[Item](fmt.Sprintf(item_url, id))
}

// GetTopStories returns up to the top 500 stories on Hacker News.
func GetTopStoriesIDs() ([]int, error) {
	return RetrieveIDs(top_stories_url)
}

// GetBestStories returns up to the best 500 stories on Hacker News.
func GetBestStoriesIDs() ([]int, error) {
	return RetrieveIDs(best_stories_url)
}

// GetNewStories returns up to the newest 500 stories on Hacker News.
func GetNewStoriesIDs() ([]int, error) {
	return RetrieveIDs(new_stories_url)
}

// GetAskStories returns up to 200 of the latest Ask stories on Hacker News.
func GetAskStoriesIDs() ([]int, error) {
	return RetrieveIDs(ask_stories_url)
}

// GetShowStories returns up to 200 of the latest Show stories on Hacker News.
func GetShowStoriesIDs() ([]int, error) {
	return RetrieveIDs(show_stories_url)
}

// GetJobStories returns up to 200 of the latest Job stories on Hacker News.
func GetJobStoriesIDs() ([]int, error) {
	return RetrieveIDs(job_stories_url)
}
