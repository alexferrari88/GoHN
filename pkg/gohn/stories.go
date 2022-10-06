package gohn

import (
	"context"
	"sort"
)

// URLs to retrieve the IDs of the top stories, best stories, new stories, ask stories, show stories and job stories.
const (
	TOP_STORIES_URL  = "topstories.json"
	BEST_STORIES_URL = "beststories.json"
	NEW_STORIES_URL  = "newstories.json"
	ASK_STORIES_URL  = "askstories.json"
	SHOW_STORIES_URL = "showstories.json"
	JOB_STORIES_URL  = "jobstories.json"
)

// StoriesService provides access to the stories endpoints of the Hacker News API.
type StoriesService service

// Story represents a story on Hacker News.
// It contains the pointer to the story item
// as Parent and a map of comments by ID as CommentsByIdMap.
type Story struct {
	Parent          *Item
	CommentsByIdMap ItemsIndex
}

// SetCommentsPosition calculates the order of the comments in a story.
// The order is calculated by traversing the Kids slice of the
// Story.Parent item recursively (N-ary tree preorder traversal).
// The order is stored in the Position field of the Item struct.
func (s *Story) SetCommentsPosition() {
	var n int
	var preorder func(int, int)
	preorder = func(id int, order int) {
		comment := s.CommentsByIdMap[id]
		comment.Position = &order
		s.CommentsByIdMap[id] = comment
		if comment.Kids != nil {
			for _, kid := range *comment.Kids {
				preorder(kid, n+1)
				n++
			}
		}
	}
	if s.Parent.Kids != nil {
		for _, kid := range *s.Parent.Kids {
			preorder(kid, n)
			n++
		}
	}
}

// GetOrderedCommentsIDs orders the comments in a Story by their Position field
// and returns a slice of comments IDs in that order.
func (s *Story) GetOrderedCommentsIDs() ([]int, error) {
	var comments []int
	for _, comment := range s.CommentsByIdMap {
		if comment.ID == nil || comment.Position == nil {
			return nil, &ErrInvalidItem{Message: "comment ID or position is nil"}
		}
		comments = append(comments, *comment.ID)
	}
	sort.Slice(comments, func(i, j int) bool {
		return *(s.CommentsByIdMap[comments[i]].Position) < *(s.CommentsByIdMap[comments[j]].Position)

	})
	return comments, nil

}

// IsTopLevelComment checks if an Item is a top level comment in a story.
func (s *Story) IsTopLevelComment(item *Item) (bool, error) {
	if s.Parent.Kids == nil {
		return false, &ErrInvalidItem{Message: "story has no kids"}
	}
	if item.Type == nil {
		return false, &ErrInvalidItem{Message: "item has no type"}
	}
	if item.Parent == nil {
		return false, &ErrInvalidItem{Message: "item has no parent"}
	}
	if *item.Type != "comment" {
		return false, &ErrInvalidItem{Message: "item is not a comment"}
	}
	for _, kid := range *s.Parent.Kids {
		if kid == *item.ID {
			return true, nil
		}
	}
	return false, nil
}

// GetTopIDs returns the IDs of up to 500 of the top stories on Hacker News.
func (s *StoriesService) GetTopIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, TOP_STORIES_URL)
}

// GetBestIDs returns the IDs of up to 500 of the best stories on Hacker News.
func (s *StoriesService) GetBestIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, BEST_STORIES_URL)
}

// GetNewIDs returns the IDs of up to 500 of the newest stories on Hacker News.
func (s *StoriesService) GetNewIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, NEW_STORIES_URL)
}

// GetAskIDs returns the IDs of up to 200 of the latest Ask stories on Hacker News.
func (s *StoriesService) GetAskIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, ASK_STORIES_URL)
}

// GetShowIDs returns the IDs of up to 200 of the latest Show stories on Hacker News.
func (s *StoriesService) GetShowIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, SHOW_STORIES_URL)
}

// GetJobIDs returns the IDs of up to 200 of the latest Job stories on Hacker News.
func (s *StoriesService) GetJobIDs(ctx context.Context) ([]*int, error) {
	return s.GetIDsFromURL(ctx, JOB_STORIES_URL)
}

// GetIDsFromURL returns a slice of IDs from a Hacker News API endpoint.
func (s *StoriesService) GetIDsFromURL(ctx context.Context, url string) ([]*int, error) {
	req, err := s.client.NewRequest("GET", url)
	if err != nil {
		return nil, err
	}
	var ids []*int
	_, err = s.client.Do(ctx, req, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
