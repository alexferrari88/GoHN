package gohn

import "sort"

// URLs to retrieve the IDs of the top stories, best stories, new stories, ask stories, show stories and job stories.
const (
	TOP_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/topstories.json"
	BEST_STORIES_URL = "https://hacker-news.firebaseio.com/v0/beststories.json"
	NEW_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/newstories.json"
	ASK_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/askstories.json"
	SHOW_STORIES_URL = "https://hacker-news.firebaseio.com/v0/showstories.json"
	JOB_STORIES_URL  = "https://hacker-news.firebaseio.com/v0/jobstories.json"
)

// Story represents a story on Hacker News.
// It contains the story item as StoryItem and a map
// of comments by ID as CommentsByIdMap.
type Story struct {
	StoryItem       Item
	CommentsByIdMap ItemsIndex
}

// SetCommentsPosition calculates the order of the comments in a story.
// The order is calculated by traversing the Kids slice of the Story item
// recursively (N-ary tree preorder traversal).
// The order is stored in the Position field of the Item struct.
func (s *Story) SetCommentsPosition() {
	var n int
	var preorder func(int, int)
	preorder = func(id int, order int) {
		comment := s.CommentsByIdMap[id]
		comment.Position = order
		s.CommentsByIdMap[id] = comment
		for _, kid := range comment.Kids {
			preorder(kid, n+1)
			n++
		}
	}
	for _, kid := range s.StoryItem.Kids {
		preorder(kid, n)
		n++
	}
}

// GetOrderedCommentsIDs orders the comments in a story by their Position field
// and returns a slice of comments IDs.
func (s *Story) GetOrderedCommentsIDs() []int {
	var comments []int
	for _, comment := range s.CommentsByIdMap {
		comments = append(comments, comment.ID)
	}
	sort.Slice(comments, func(i, j int) bool {
		return s.CommentsByIdMap[comments[i]].Position < s.CommentsByIdMap[comments[j]].Position
	})
	return comments

}

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
