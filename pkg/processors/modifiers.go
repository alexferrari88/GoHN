package processors

import (
	"html"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

// UnescapeHTML unescapes HTML entities in the text of the item
func UnescapeHTML() gohn.ItemProcessor {
	return func(item *gohn.Item) error {
		item.Text = html.UnescapeString(item.Text)
		return nil
	}
}
