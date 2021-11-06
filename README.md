# hnscraper

[![Go Reference](https://pkg.go.dev/badge/github.com/thetallpaul/hnscraper.svg)](https://pkg.go.dev/github.com/thetallpaul/hnscraper)

Web scraper for HackerNews.
While HackerNews has a [fantastic API](https://github.com/HackerNews/API), maybe you'd prefer to scrape the pages directly instead?

Using `hnscraper` is simple. If you want to request a single page, use `ScrapePage()`:

```go
package main

import (
  "fmt"

  "github.com/thetallpaul/hnscaper"
)

func main() {
  pageTwo := hnscraper.ScapePage(2)

  // Prints the first title on the second page
  fmt.Println(pageTwo.Posts[0].Title)
}
```

If you want to select an inclusive range of pages, use `ScrapeMultPages()`:

```go
pages := hnscaper.ScrapeMultPages(2,4)

// Prints the number of votes on every post for pages 2-4
for _, page := range pages {
  for _, post := range page.Posts {
    fmt.Printf("Page: %d, Rank: %d, Votes: %d\n", page.Num, post.Rank, post.Votes)
  }
}
```
