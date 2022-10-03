package processorstest

import (
	"context"
	"html"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
	"github.com/alexferrari88/gohn/pkg/processors"
)

func TestUnescapeHTML(t *testing.T) {
	mockItemToTestID := 8
	mockItem, _, mockClient, err := setup([]gohn.Item{
		{
			ID:   mockItemToTestID,
			Text: html.EscapeString(`This is an <a href="https://www.example.com">example</a>`),
		},
	})

	if err != nil {
		t.Fatalf("error setting up test: %v", err)
	}

	client := gohn.NewClient(context.Background(), mockClient)
	items := client.RetrieveKidsItems(mockItem, processors.UnescapeHTML())

	if len(items) != 7 {
		t.Fatalf("expected 7 items, got %v", len(items))
	}

	if items[mockItemToTestID].Text != `This is an <a href="https://www.example.com">example</a>` {
		t.Errorf("expected unescaped text, got %v", items[mockItemToTestID].Text)
	}

}
