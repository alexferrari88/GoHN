package itemprocessors

import (
	"html"

	"github.com/alexferrari88/gohn"
)

// UnescapeHTML unescapes HTML entities in the text of the item
func UnescapeHTML() gohn.ItemProcessor {
	return func(item *gohn.Item) error {
		item.Text = html.UnescapeString(item.Text)
		return nil
	}
}
