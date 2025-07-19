package processors

import (
	"sync"
	"testing"

	"github.com/alexferrari88/gohn/pkg/gohn"
)

// BenchmarkFilterOutWords_NilItem tests performance with nil items (early return)
func BenchmarkFilterOutWords_NilItem(b *testing.B) {
	processor := FilterOutWords([]string{"test"}, false)
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(nil, &wg)
	}
}

// BenchmarkFilterOutWords_NoText tests performance with items that have no text (early return)
func BenchmarkFilterOutWords_NoText(b *testing.B) {
	processor := FilterOutWords([]string{"test"}, false)
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID: intPtr(1),
		// No Text field set
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
	}
}

// BenchmarkFilterOutWords_WithText tests performance with items that have text (full processing)
func BenchmarkFilterOutWords_WithText(b *testing.B) {
	processor := FilterOutWords([]string{"test"}, false)
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID:   intPtr(1),
		Text: stringPtr("This is some sample text without the target word"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
		wg.Wait() // Wait for the WaitGroup to complete
	}
}

// BenchmarkFilterOutDeleted_NilItem tests performance with nil items (early return)
func BenchmarkFilterOutDeleted_NilItem(b *testing.B) {
	processor := FilterOutDeleted()
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(nil, &wg)
	}
}

// BenchmarkFilterOutDeleted_ValidItem tests performance with valid items (full processing)
func BenchmarkFilterOutDeleted_ValidItem(b *testing.B) {
	processor := FilterOutDeleted()
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID:      intPtr(1),
		Deleted: boolPtr(false),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
		wg.Wait() // Wait for the WaitGroup to complete
	}
}

// BenchmarkFilterOutUsers_NilItem tests performance with nil items (early return)
func BenchmarkFilterOutUsers_NilItem(b *testing.B) {
	processor := FilterOutUsers([]string{"testuser"})
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(nil, &wg)
	}
}

// BenchmarkFilterOutUsers_NoAuthor tests performance with items that have no author (early return)
func BenchmarkFilterOutUsers_NoAuthor(b *testing.B) {
	processor := FilterOutUsers([]string{"testuser"})
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID: intPtr(1),
		// No By field set
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
	}
}

// BenchmarkFilterOutUsers_WithAuthor tests performance with items that have an author (full processing)
func BenchmarkFilterOutUsers_WithAuthor(b *testing.B) {
	processor := FilterOutUsers([]string{"testuser"})
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID: intPtr(1),
		By: stringPtr("normaluser"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
		wg.Wait() // Wait for the WaitGroup to complete
	}
}

// BenchmarkUnescapeHTML_NilItem tests performance with nil items (early return)
func BenchmarkUnescapeHTML_NilItem(b *testing.B) {
	processor := UnescapeHTML()
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(nil, &wg)
	}
}

// BenchmarkUnescapeHTML_NoText tests performance with items that have no text (early return)
func BenchmarkUnescapeHTML_NoText(b *testing.B) {
	processor := UnescapeHTML()
	var wg sync.WaitGroup
	
	item := &gohn.Item{
		ID: intPtr(1),
		// No Text field set
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor(item, &wg)
	}
}

// BenchmarkUnescapeHTML_WithText tests performance with items that have text (full processing)
func BenchmarkUnescapeHTML_WithText(b *testing.B) {
	processor := UnescapeHTML()
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a copy of the text for each iteration to avoid modifying the original
		textCopy := "This &amp; that &lt;tag&gt;"
		item := &gohn.Item{
			ID:   intPtr(1),
			Text: &textCopy,
		}
		processor(item, &wg)
		wg.Wait() // Wait for the WaitGroup to complete
	}
}

// Helper functions to create pointers
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}