package main

import (
	"fmt"
	"html"

	"github.com/alexferrari88/gohn"
)

func main() {
	topStoriesIds, _ := gohn.GetTopStoriesIDs()
	story, _ := gohn.GetItem(topStoriesIds[0])
	commentsMap := story.RetrieveKidsItems()
	fmt.Printf("Comments found: %d\n", len(commentsMap))
	fmt.Println()
	// preorder traversal of the n-ary tree story.Kids
	// to print the comments
	var preorder func(int)
	preorder = func(id int) {
		comment := commentsMap[id]
		fmt.Println(html.UnescapeString(comment.Text))
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
