package processors

import (
	"html"
	"sync"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

// UnescapeHTML unescapes HTML entities in the text of the item
func UnescapeHTML() gohn.ItemProcessor {
	return func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		// Early return for nil items - no WaitGroup operations needed
		if item == nil {
			return false, nil
		}
		
		// Early return if item has no text - no WaitGroup operations needed
		if item.Text == nil {
			return false, nil
		}
		
		// Only use WaitGroup when we actually need to process the text
		wg.Add(1)
		defer wg.Done()
		
		*item.Text = html.UnescapeString(*item.Text)
		return false, nil
	}
}
