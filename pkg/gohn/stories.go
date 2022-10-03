// Package gohn is a wrapper for the Hacker News API which uses goroutines.
package gohn

const (
	top_stories_url  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	best_stories_url = "https://hacker-news.firebaseio.com/v0/beststories.json"
	new_stories_url  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	ask_stories_url  = "https://hacker-news.firebaseio.com/v0/askstories.json"
	show_stories_url = "https://hacker-news.firebaseio.com/v0/showstories.json"
	job_stories_url  = "https://hacker-news.firebaseio.com/v0/jobstories.json"
	item_url         = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

// GetTopStories returns the IDs of up to 500 of the top stories on Hacker News.
func GetTopStoriesIDs() ([]int, error) {
	return RetrieveIDs(top_stories_url)
}

// GetBestStories returns the IDs of up to 500 of the best stories on Hacker News.
func GetBestStoriesIDs() ([]int, error) {
	return RetrieveIDs(best_stories_url)
}

// GetNewStories returns the IDs of up to 500 of the newest stories on Hacker News.
func GetNewStoriesIDs() ([]int, error) {
	return RetrieveIDs(new_stories_url)
}

// GetAskStories returns the IDs of up to 200 of the latest Ask stories on Hacker News.
func GetAskStoriesIDs() ([]int, error) {
	return RetrieveIDs(ask_stories_url)
}

// GetShowStories returns the IDs of up to 200 of the latest Show stories on Hacker News.
func GetShowStoriesIDs() ([]int, error) {
	return RetrieveIDs(show_stories_url)
}

// GetJobStories returns the IDs of up to 200 of the latest Job stories on Hacker News.
func GetJobStoriesIDs() ([]int, error) {
	return RetrieveIDs(job_stories_url)
}
