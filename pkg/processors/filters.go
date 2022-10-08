package processors

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

// FilterOutWords filters items that contain the given words in the title or text.
// The argument title is a boolean that indicates if the filter
// should be applied to the title and not the text.
func FilterOutWords(words []string, title bool) gohn.ItemProcessor {
	return func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		if item == nil {
			return true, nil
		}
		wg.Add(1)
		defer wg.Done()
		for _, word := range words {
			if title && item.Title != nil {
				if strings.Contains(strings.ToLower(*item.Title), strings.ToLower(word)) {
					return false, fmt.Errorf("word %s found in title", word)
				}
			} else {
				if item.Text == nil {
					return false, nil
				}
				if strings.Contains(strings.ToLower(*item.Text), strings.ToLower(word)) {
					return false, fmt.Errorf("Word found")
				}
			}
		}
		return false, nil
	}
}

// FilterOutDeleted filters deleted items
func FilterOutDeleted() gohn.ItemProcessor {
	return func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		if item == nil {
			return true, nil
		}
		wg.Add(1)
		defer wg.Done()
		if item.Deleted != nil && *item.Deleted {
			return false, fmt.Errorf("Deleted item found")
		}
		return false, nil
	}
}

// FilterOutUsers filters items that are not from the given user
func FilterOutUsers(users []string) gohn.ItemProcessor {
	return func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		if item == nil {
			return true, nil
		}
		if item.By == nil {
			return false, nil
		}
		wg.Add(1)
		defer wg.Done()
		for _, user := range users {
			if *item.By == user {
				return false, fmt.Errorf("User found")
			}
		}
		return false, nil
	}
}
