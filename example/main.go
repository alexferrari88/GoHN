package main

import (
	"fmt"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

func main() {
	// Instantiate a new client
	client := gohn.NewDefaultClient()

	// Get the top 500 stories' IDs
	topStoriesIds, _ := client.GetTopStoriesIDs()

	// Retrieve the first one
	story, _ := client.GetItem(topStoriesIds[0])

	// Retrieve all the comments for that story
	// UnescapeHTML is applied to each retrieved item to unescape HTML characters
	commentsMap := client.RetrieveKidsItems(story, processors.UnescapeHTML())

	fmt.Printf("Comments found: %d\n", len(commentsMap))
	fmt.Println()

	// preorder traversal of the n-ary tree story.Kids
	// to print the comments
	var preorder func(int)
	preorder = func(id int) {
		comment := commentsMap[id]
		fmt.Println(comment.Text)
		for _, kid := range comment.Kids {
			fmt.Print("\t")
			preorder(kid)
		}
	}
	for _, kid := range story.Kids {
		fmt.Print("> ")
		preorder(kid)
		fmt.Println()
	}

}
