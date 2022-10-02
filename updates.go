package gohn

// Update represents a items and profiles changes from the Hacker News API.
type Update struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

// updates_url is the URL for the updates endpoint.
const (
	updates_url = "https://hacker-news.firebaseio.com/v0/updates.json"
)

// GetUpdates returns a slice of IDs for the given URL.
func GetUpdates() (Update, error) {
	return retrieveFromURL[Update](updates_url)
}
