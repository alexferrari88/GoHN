# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoHN is a Go wrapper for the Hacker News API, inspired by go-github. It provides a simple interface to fetch stories, comments, and user data from Hacker News with features like concurrent comment retrieval and filtering capabilities.

## Key Architecture

- **Core Client**: `pkg/gohn/hn.go` - Main HTTP client with rate limiting (1 req/sec)
- **Services**: Organized by data type (Items, Stories, Users, Updates)
  - `pkg/gohn/items.go` - Fetch individual items and comments with goroutines
  - `pkg/gohn/stories.go` - Get story lists (top, new, best, ask, show, job)
  - `pkg/gohn/users.go` - User profile operations
  - `pkg/gohn/updates.go` - Recent updates tracking
- **Processors**: `pkg/processors/` - Filters and modifiers for retrieved data
- **Story Structure**: `Story` type combines parent item with ordered comments map

## Common Commands

### Build and Test
```bash
go build ./...                    # Build all packages
go test ./...                     # Run all tests
go test ./test/gohn              # Test core functionality
go test ./test/processors        # Test filters/processors
go test -v ./...                 # Verbose test output
```

### Development
```bash
go run example/main.go           # Run example application
go mod tidy                      # Clean up dependencies
go vet ./...                     # Static analysis
go fmt ./...                     # Format code
```

### Performance Testing
```bash
go test -bench=. ./...           # Run benchmarks
go test -cpuprofile=cpu.prof ./test/gohn  # CPU profiling
```

## Key Patterns

- All API methods take `context.Context` as first parameter
- Rate limiting is built into the client (1 request per second)
- Use `FetchAllDescendants()` for concurrent comment retrieval
- Apply processors like `processors.UnescapeHTML()` for data cleaning
- Stories organize comments using `SetCommentsPosition()` and `GetOrderedCommentsIDs()`

## Testing Structure

Tests are in `test/` directory organized by package. Use `test/setup/setup.go` for test utilities and shared setup code.