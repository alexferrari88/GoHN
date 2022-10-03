package main

import (
	"fmt"

	"github.com/alexferrari88/gohn"
	itemprocessors "github.com/alexferrari88/gohn/item-processors"
)

func main() {
	topStoriesIds, _ := gohn.GetTopStoriesIDs()
	story, _ := gohn.GetItem(topStoriesIds[0])
	// UnescapeHTML is applied to each retrieved item to unescape HTML characters
	commentsMap := story.RetrieveKidsItems(itemprocessors.UnescapeHTML())
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
