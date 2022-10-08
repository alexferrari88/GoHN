package processors

import (
	"html"
	"sync"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

// UnescapeHTML unescapes HTML entities in the text of the item
func UnescapeHTML() gohn.ItemProcessor {
	return func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		if item == nil {
			return false, nil
		}
		if item.Text == nil {
			return false, nil
		}
		wg.Add(1)
		defer wg.Done()
		*item.Text = html.UnescapeString(*item.Text)
		return false, nil
	}
}
