package processors

import (
	"sync"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

func TestUnescapeHTML(t *testing.T) {
	expectedText := `This is an <a href="https://www.example.com">example</a>`
	var wg sync.WaitGroup
	id := 1
	i := &gohn.Item{ID: &id, Text: &expectedText}

	f := UnescapeHTML()
	err := f(i, &wg)

	if err != nil {
		t.Fatalf("unexpected error unescaping HTML: %v", err)
	}

	if *i.Text != expectedText {
		t.Fatalf("expected text to be %v, got %v", expectedText, *i.Text)
	}
}
