package hnscraper

import (
	"testing"
)

func TestScrapePageFail(t *testing.T) {
	_, err := ScrapePage(0)

	if err == nil {
		t.Error("accepted invalid page number")
		return
	}
}

func TestScrapePageSuccess(t *testing.T) {
	result, err := ScrapePage(1)

	if err != nil {
		t.Error("error: ", err)
		return
	}

	numPosts := len(result.Posts)
	if numPosts < 20 {
		t.Error("returned insufficient number of posts (", numPosts, ") for frontpage")
	}
}

func TestScrapeMultPagesFail(t *testing.T) {
	_, err := ScrapeMultPages(-1, 2)

	if err == nil {
		t.Error("accepted invalid page range")
		return
	}
}

func TestScrapeMultPagesSuccess(t *testing.T) {
	result, err := ScrapeMultPages(1, 3)

	if err != nil {
		t.Error("error: ", err)
		return
	}

	numPages := len(result)
	if numPages != 3 {
		t.Error("returned ", numPages, " pages instead of 3")
	}
}
