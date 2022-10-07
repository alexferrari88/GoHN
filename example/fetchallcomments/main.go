package main

import (
	"context"
	"fmt"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

// Retrieve all the comments for that story
func fetchComments() ([]string, error) {

	// Instantiate a new client to retrieve data from the Hacker News API
	hn := gohn.NewClient(nil)

	// Use background context
	ctx := context.Background()

	// Get the top 500 stories' IDs
	topStoriesIds, _ := hn.Stories.GetTopIDs(ctx)

	var story *gohn.Item

	var comments []string

	// Retrieve the details of the first one
	if len(topStoriesIds) > 0 && topStoriesIds[0] != nil {
		story, _ = hn.Items.Get(ctx, *topStoriesIds[0])
	}

	if story == nil {
		panic("No story found")
	}

	commentsMap, err := hn.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())

	if err != nil {
		panic(err)
	}
	if len(commentsMap) == 0 {
		fmt.Println("No comments found")
	}

	fmt.Printf("Comments found: %d\n", len(commentsMap))

	// Create a Story struct to hold the story and its comments
	storyWithComments := gohn.Story{
		Parent:          story,
		CommentsByIdMap: commentsMap,
	}

	// Calculate the position of each comment in the story
	storyWithComments.SetCommentsPosition()

	// Get an ordered list of comments' IDs (ordered by position)
	orderedIDs, err := storyWithComments.GetOrderedCommentsIDs()

	if err != nil {
		panic(err)
	}

	// Print the comments
	for _, id := range orderedIDs {
		comment := commentsMap[id]
		if comment.Text != nil {
			comments = append(comments, *comment.Text)
		}
	}
	return comments, err

}

func main() {
	comments, err := fetchComments()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i := range comments {
		fmt.Printf("\n%v: %s \n\n", i+1, comments[i])
	}
}
