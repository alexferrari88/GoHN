package itemprocessors

import (
	"fmt"

	"github.com/alexferrari88/gohn"
)

// FilterOutWords filters items that contain the given words in the title or text.
// The argument title is a boolean that indicates if the filter
// should be applied to the title and not the text.
func FilterOutWords(words []string, title bool) gohn.ItemProcessor {
	return func(item *gohn.Item) error {
		for _, word := range words {
			if title {
				if item.Title == word {
					return fmt.Errorf("word %s found in title", word)
				}
			} else {
				if item.Text == word {
					return fmt.Errorf("Word found")
				}
			}
		}
		return nil
	}
}

// FilterOutDeleted filters deleted items
func FilterOutDeleted() gohn.ItemProcessor {
	return func(item *gohn.Item) error {
		if item.Deleted {
			return fmt.Errorf("Deleted item found")
		}
		return nil
	}
}

// FilterOutUsers filters items that are not from the given user
func FilterOutUsers(users []string) gohn.ItemProcessor {
	return func(item *gohn.Item) error {
		for _, user := range users {
			if item.By == user {
				return fmt.Errorf("User found")
			}
		}
		return nil
	}
}
