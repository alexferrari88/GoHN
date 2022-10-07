package main

import (
	"context"
	"fmt"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

func main() {
	// Instantiate a new client to retrieve data from the Hacker News API
	hn := gohn.NewClient(nil)

	// Use background context
	ctx := context.Background()

	// Get the top 500 stories' IDs
	topStoriesIds, _ := hn.Stories.GetTopIDs(ctx)

	var story *gohn.Item
	// Retrieve the details of the first one
	if len(topStoriesIds) > 0 && topStoriesIds[0] != nil {
		story, _ = hn.Items.Get(ctx, *topStoriesIds[0])
	}

	if story == nil {
		panic("No story found")
	}

	// Print the story's title
	fmt.Println("Title:", *story.Title)

	// Print the story's author
	fmt.Println("Author:", *story.By)

	// Print the story's score
	fmt.Println("Score:", *story.Score)

	// Print the story's URL
	fmt.Println("URL:", *story.URL)

	fmt.Println()
	fmt.Println()

	if story.Kids == nil {
		fmt.Println("No comments found")
		return
	}

	// Retrieve all the comments for that story
	// UnescapeHTML is applied to each retrieved item to unescape HTML characters

	commentsMap, err := hn.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())

	if err != nil {
		panic(err)
	}
	if len(commentsMap) == 0 {
		fmt.Println("No comments found")
		return
	}

	fmt.Printf("Comments found: %d\n", len(commentsMap))
	fmt.Println()

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
			fmt.Println(*comment.Text)
			fmt.Println()
		}
	}

}
