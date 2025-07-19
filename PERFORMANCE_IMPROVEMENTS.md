# GoHN Performance Technical Details

This document provides technical details about the performance optimizations implemented in GoHN for developers who need to understand the implementation or optimize for specific use cases.

## Benchmark Results

**Test case**: HN Thread 34591291 (271 comments)

| Metric | Original | Optimized | Improvement |
|--------|----------|-----------|-------------|
| **Time (10 comments)** | 9.003s | 0.090s | **100x faster** |
| **Memory** | 125,112 bytes | 97,308 bytes | **22% reduction** |
| **Allocations** | 1,199 | 1,081 | **10% reduction** |
| **Rate limit** | 1 req/sec | Configurable | **Up to 10x throughput** |

## Technical Optimizations

### 1. Configurable Rate Limiting
- `NewClientWithOptions()` allows custom rate limits and burst sizes
- Default: 1 req/sec → Configurable up to 10+ req/sec
- HTTP transport optimization with connection pooling

### 2. Worker Pool Architecture  
- `FetchAllDescendantsOptimized()` uses fixed worker pools instead of unlimited goroutines
- Automatic worker count calculation: `min(max(descendants/20, 5), 50)`
- Reduces goroutine overhead and improves resource utilization

### 3. Memory Pre-allocation
- Channel buffers pre-sized based on `descendants` count
- Map pre-allocation: `make(ItemsIndex, expectedCount)`
- Buffer size capping to prevent excessive memory usage

### 4. Concurrent Architecture Improvements
- Separated work distribution from result collection
- Optimized channel communication patterns
- Better error handling and context cancellation

## Configuration Guidelines

### Worker Count Optimization
```go
// Auto-calculation (recommended)
workerCount := 0  // Library determines optimal count

// Manual calculation for custom scenarios
workerCount := min(max(expectedComments/20, 5), 50)
```

### Rate Limit Recommendations
| Use Case | Rate Limit | Burst Size | Comments |
|----------|------------|------------|----------|
| Development/Testing | 100ms | 20 | 10 req/sec (use with caution) |
| Production | 200ms | 10 | 5 req/sec (recommended) |
| Conservative | 1000ms | 1 | 1 req/sec (default) |

### Memory Optimization
- Channel buffers capped at 1000 items to prevent excessive memory usage
- Map pre-allocation based on `descendants` field when available
- Fallback to `len(kids) * 2` estimation when `descendants` unavailable

## Performance Characteristics

### Bottleneck Analysis
- **Primary bottleneck**: API rate limiting (not CPU or memory)
- **Secondary**: HTTP connection overhead (mitigated by connection pooling)
- **Minimal impact**: Goroutine creation overhead (optimized with worker pools)

### Scaling Behavior
- **Linear scaling** with rate limit increases
- **Diminishing returns** beyond 10 workers for typical comment trees
- **Memory usage**: ~350 bytes per comment + HTTP response buffers