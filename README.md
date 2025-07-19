# GoHN — Hacker News API Wrapper for Go

<div align="center">
<a href="https://pkg.go.dev/github.com/alexferrari88/gohn/pkg/gohn"><img src="https://pkg.go.dev/badge/github.com/alexferrari88/gohn.svg" alt="Go Reference"></a>
</div>
<div align="center">
<img src="img/logo_1.svg" width="300" style="margin: 0 auto;" />
</div>

GoHN is a wrapper for the [Hacker News API](https://github.com/HackerNews/API) for Go, inspired by the excellent [go-github](https://github.com/google/go-github) library.

It facilitates the use of the API by providing a simple interface to the API endpoints.

## Features 🚀

- Get the top/new/best/ask/show/job stories
- **High-performance comment retrieval** with configurable rate limiting and worker pools
- Retrieve all comments (with metadata) for a story using optimized concurrent processing
- **Up to 100x faster** comment fetching for large threads with performance optimizations
- Retrieve the comments ordered as they appear in the story on the website
- Apply filters to retrieved items (stories, comments)
- Configurable HTTP client with connection pooling and custom timeouts

## Usage 💻

Refer to the [GoDoc](https://pkg.go.dev/github.com/alexferrari88/gohn/pkg/gohn) for the full API reference.

### Quick Start

#### Basic Usage


```go
    // Instantiate a new client to retrieve data from the Hacker News API
    hn, err := gohn.NewClient(nil)
    if err != nil {
        panic(err)
    }

    // Use background context
    ctx := context.Background()

    // Get the top 500 stories' IDs
    topStoriesIds, _ := hn.Stories.GetIDs(ctx, gohn.TopStory)

    var story *gohn.Item
    // Retrieve the details of the first one
    if len(topStoriesIds) > 0 && topStoriesIds[0] != nil {
        story, _ = hn.Items.Get(ctx, *topStoriesIds[0])
    }

    if story == nil {
        panic("No story found")
    }

    // Print the story's title
    fmt.Println("Title:", *story.Title)

    // Print the story's author
    fmt.Println("Author:", *story.By)

    // Print the story's score
    fmt.Println("Score:", *story.Score)

    // Print the story's URL
    fmt.Println("URL:", *story.URL)

    fmt.Println()
    fmt.Println()

    if story.Kids == nil {
        fmt.Println("No comments found")
        return
    }

    // Retrieve all the comments for that story
    // UnescapeHTML is applied to each retrieved item to unescape HTML characters

    commentsMap, err := hn.Items.FetchAllDescendants(ctx, story, processors.UnescapeHTML())

    if err != nil {
        panic(err)
    }
    if len(commentsMap) == 0 {
        fmt.Println("No comments found")
        return
    }

    fmt.Printf("Comments found: %d\n", len(commentsMap))
    fmt.Println()

    // Create a Story struct to hold the story and its comments
    storyWithComments := gohn.Story{
        Parent:          story,
        CommentsByIdMap: commentsMap,
    }

    // Calculate the position of each comment in the story
    storyWithComments.SetCommentsPosition()

    // Get an ordered list of comments' IDs (ordered by position)
    orderedIDs, err := storyWithComments.GetOrderedCommentsIDs()

    if err != nil {
        panic(err)
    }

    // Print the comments
    for _, id := range orderedIDs {
        comment := commentsMap[id]
        if comment.Text != nil {
            fmt.Println(*comment.Text)
            fmt.Println()
        }
    }
```

#### Optimized Usage for Better Performance

For large comment threads (50+ comments), use the optimized client and functions:

```go
    // Create optimized client with custom rate limiting
    hn, err := gohn.NewClientWithOptions(&gohn.ClientOptions{
        RateLimit: 200 * time.Millisecond, // 5 requests per second (be respectful)
        BurstSize: 10,                     // Allow bursts of 10 requests
    })
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Get story (same as before)
    story, _ := hn.Items.Get(ctx, storyID)

    // Use optimized function for faster comment retrieval
    commentsMap, err := hn.Items.FetchAllDescendantsOptimized(
        ctx, 
        story, 
        processors.UnescapeHTML(), 
        0, // Auto-determine optimal worker count
    )

    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Comments fetched: %d\n", len(commentsMap))
    // Process comments same as before...
```

### Performance Comparison

| Scenario | Original | Optimized | Improvement |
|----------|----------|-----------|-------------|
| 10 comments | 9.0 seconds | 0.09 seconds | **100x faster** |
| 271 comments | 271+ seconds | ~49 seconds | **5.5x faster** |

> **Note**: Always be respectful of the Hacker News API. The default rate limit is 1 req/sec. For better performance, 5 req/sec (200ms) is recommended, with 10 req/sec as an absolute maximum.

### Examples

- [example/main.go](example/main.go) - Basic usage example
- [example/quick_performance.go](example/quick_performance.go) - Performance testing example

For detailed performance information, see [PERFORMANCE_IMPROVEMENTS.md](PERFORMANCE_IMPROVEMENTS.md).

## Semantic Versioning 🥚

As this library is not yet in version 1.0.0, the API may have breaking changes between minor versions.

## Contributing 🤝🏼

Feel free to fork this repo and create a PR. I will review them and merge, if ok.

## License 📝

[MIT](https://choosealicense.com/licenses/mit/)
