package main

import "github.com/alexferrari88/gohn"

func main() {
	// Get the top 10 stories.
	ids, _ := gohn.GetTopStoriesIDs()
	for _, id := range ids[:10] {
		item, _ := gohn.GetItem(id)
		println(item.Title)
	}
}
