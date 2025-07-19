package gohn

import (
	"context"
	"fmt"
	"sort"
)

// StoryType represents the type of a story.
type StoryType string

const (
	TopStory  StoryType = "topstories"
	BestStory StoryType = "beststories"
	NewStory  StoryType = "newstories"
	AskStory  StoryType = "askstories"
	ShowStory StoryType = "showstories"
	JobStory  StoryType = "jobstories"
)

// StoriesService provides access to the stories endpoints of the Hacker News API.
type StoriesService service

// ItemsIndex is a map of Items indexed by their ID.
type ItemsIndex map[int]*Item

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
		if comment, ok := s.CommentsByIdMap[id]; ok {
			comment.Position = &order
			s.CommentsByIdMap[id] = comment
			if comment.Kids != nil {
				for _, kid := range *comment.Kids {
					preorder(kid, n+1)
					n++
				}
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
			return nil, &InvalidItemError{Message: "comment ID or position is nil"}
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
		return false, &InvalidItemError{Message: "story has no kids"}
	}
	if item.Type == nil {
		return false, &InvalidItemError{Message: "item has no type"}
	}
	if item.Parent == nil {
		return false, &InvalidItemError{Message: "item has no parent"}
	}
	if *item.Type != "comment" {
		return false, &InvalidItemError{Message: "item is not a comment"}
	}
	for _, kid := range *s.Parent.Kids {
		if kid == *item.ID {
			return true, nil
		}
	}
	return false, nil
}

// GetIDs returns the IDs of stories of a given type.
func (s *StoriesService) GetIDs(ctx context.Context, storyType StoryType) ([]*int, error) {
	url := fmt.Sprintf("%s.json", storyType)
	return s.getIDsFromURL(ctx, url)
}

// getIDsFromURL returns a slice of IDs from a Hacker News API endpoint.
func (s *StoriesService) getIDsFromURL(ctx context.Context, url string) ([]*int, error) {
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