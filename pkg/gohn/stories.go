// Package gohn is a wrapper for the Hacker News API which uses goroutines.
package gohn

const (
	TOP_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	BEST_STORIES_URL = "https://hacker-news.firebaseio.com/v0/beststories.json"
	NEW_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	ASK_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/askstories.json"
	SHOW_STORIES_URL = "https://hacker-news.firebaseio.com/v0/showstories.json"
	JOB_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/jobstories.json"
)

// GetTopStories returns the IDs of up to 500 of the top stories on Hacker News.
func (c client) GetTopStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(TOP_STORIES_URL)
}

// GetBestStories returns the IDs of up to 500 of the best stories on Hacker News.
func (c client) GetBestStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(BEST_STORIES_URL)
}

// GetNewStories returns the IDs of up to 500 of the newest stories on Hacker News.
func (c client) GetNewStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(NEW_STORIES_URL)
}

// GetAskStories returns the IDs of up to 200 of the latest Ask stories on Hacker News.
func (c client) GetAskStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(ASK_STORIES_URL)
}

// GetShowStories returns the IDs of up to 200 of the latest Show stories on Hacker News.
func (c client) GetShowStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(SHOW_STORIES_URL)
}

// GetJobStories returns the IDs of up to 200 of the latest Job stories on Hacker News.
func (c client) GetJobStoriesIDs() ([]int, error) {
	return c.RetrieveIDs(JOB_STORIES_URL)
}
