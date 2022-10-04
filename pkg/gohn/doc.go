/*
Package gohn is a tiny wrapper for the Hacker News API: https://github.com/HackerNews/API

	  Example usage:
	  	package main

	  	import (
	  		"fmt"

	  		"github.com/alexferrari88/gohn/pkg/gohn"
	  		"github.com/alexferrari88/gohn/pkg/processors"
	  	)

	  	func main() {

			// Instantiate a new client to retrieve data from the Hacker News API
			hn := gohn.NewDefaultClient()

			// Get the top 500 stories' IDs
			topStoriesIds, _ := hn.GetTopStoriesIDs()

			// Retrieve the details of the first one
			story, _ := hn.GetItem(topStoriesIds[0])

			// Print the story's title
			fmt.Println("Title:", story.Title)

			// Print the story's author
			fmt.Println("Author:", story.By)

			// Print the story's score
			fmt.Println("Score:", story.Score)

			// Print the story's URL
			fmt.Println("URL:", story.URL)

			fmt.Println()
			fmt.Println()

			// Retrieve all the comments for that story
			// UnescapeHTML is applied to each retrieved item to unescape HTML characters
			commentsMap := hn.RetrieveKidsItems(story, processors.UnescapeHTML())

			fmt.Printf("Comments found: %d\n", len(commentsMap))
			fmt.Println()

			// Create a Story struct to hold the story and its comments

				storyWithComments := gohn.Story{
					StoryItem:       story,
					CommentsByIdMap: commentsMap,
				}

			// Calculate the position of each comment in the story
			storyWithComments.SetCommentsPosition()

			// Get an ordered list of comments' IDs (ordered by position)
			orderedIDs := storyWithComments.GetOrderedCommentsIDs()

			// Print the comments

				for _, id := range orderedIDs {
					comment := commentsMap[id]
					fmt.Println(comment.Text)
					fmt.Println()
				}

}
*/
package gohn
