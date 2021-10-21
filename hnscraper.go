package hnscraper

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Post struct {
	Rank        int
	Title       string
	Score       int
	By          string
	URL         string
	NumComments int
	TimePosted  time.Time
}

type Page struct {
	Posts     []Post
	Num       int
	Retrieved time.Time
}

const hackernewsURL = "https://news.ycombinator.com/news?p="

func ScrapePages(startPage, endPage int) ([]Page, error) {
	var pages []Page

	if startPage < 1 || endPage < 1 {
		return pages, errors.New("page numbers must be positive integers")
	} else if startPage > endPage {
		return pages, errors.New(
			"starting page number cannot be larger than ending page number")
	}

	for i := startPage; i <= endPage; i++ {
		page, err := ScrapePage(i)
		if err != nil {
			return pages, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func ScrapePage(pageNum int) (Page, error) {
	var page Page
	var posts []Post

	if pageNum < 1 {
		return page, errors.New("page number must be a positive integer")
	}

	retrievedTime := time.Now()
	doc, err := htmlquery.LoadURL(hackernewsURL + strconv.Itoa(pageNum))
	if err != nil {
		return page, err
	}

	listNodes := htmlquery.Find(doc, "//table[contains(@class, 'itemlist')]/tbody/tr")

	for i := 0; i < len(listNodes)-2; i += 3 {
		subtext := htmlquery.FindOne(listNodes[i+1], "/td[contains(@class, 'subtext')]")
		post, err := getPost(listNodes[i], subtext)
		if err != nil {
			return page, err
		}

		posts = append(posts, post)
	}

	page = Page{Posts: posts, Num: pageNum, Retrieved: retrievedTime}
	return page, nil
}

func getPost(titleNode, subtextNode *html.Node) (Post, error) {
	var post Post

	title, err := getTitle(titleNode)
	if err != nil {
		return post, err
	}

	rank, err := getRank(titleNode)
	if err != nil {
		return post, err
	}

	url, err := getURL(titleNode)
	if err != nil {
		return post, err
	}

	author, err := getAuthor(subtextNode)
	if err != nil {
		return post, err
	}

	points, err := getPoints(subtextNode)
	if err != nil {
		return post, err
	}

	numComments, err := getNumComments(subtextNode)
	if err != nil {
		return post, err
	}

	timePosted, err := getTimePosted(subtextNode)
	if err != nil {
		return post, err
	}

	post = Post{
		Title:       title,
		Score:       points,
		Rank:        rank,
		By:          author,
		URL:         url,
		NumComments: numComments,
		TimePosted:  timePosted,
	}

	return post, nil
}

const errorMsg = "could not process: page formatted unexpectedly"

func getTitle(node *html.Node) (string, error) {
	title := ""
	titleQuery := htmlquery.Find(node, "/td/a")
	if len(titleQuery) != 1 {
		return title, errors.New(errorMsg)
	}
	title = htmlquery.InnerText(titleQuery[0])

	return title, nil
}

func getRank(node *html.Node) (int, error) {
	rank := 0
	rankQuery := htmlquery.Find(node, "/td/span[contains(@class, 'rank')]")
	if len(rankQuery) != 1 {
		return rank, errors.New(errorMsg)
	}
	rankStr := htmlquery.InnerText(rankQuery[0])
	rank, err := strconv.Atoi(rankStr[:len(rankStr)-1])
	if err != nil {
		return rank, err
	}

	return rank, nil
}

func getURL(node *html.Node) (string, error) {
	url := ""
	urlQuery := htmlquery.Find(node, "/td/a[contains(@class, 'titlelink')]")
	if len(urlQuery) != 1 {
		return url, errors.New(errorMsg)
	}
	url = htmlquery.SelectAttr(urlQuery[0], "href")

	return url, nil
}

func getAuthor(node *html.Node) (string, error) {
	author := ""
	authorQuery := htmlquery.Find(node, "/a[contains(@class, 'hnuser')]")
	if len(authorQuery) == 1 {
		author = htmlquery.InnerText(authorQuery[0])
	}

	return author, nil
}

func getPoints(node *html.Node) (int, error) {
	points := 0
	pointsQuery := htmlquery.Find(node, "/span[contains(@class, 'score')]")
	if len(pointsQuery) != 1 {
		return points, errors.New(errorMsg)
	}
	pointsStr := htmlquery.InnerText(pointsQuery[0])
	points, err := strconv.Atoi(
		strings.TrimSpace(strings.ReplaceAll(pointsStr, "points", "")))
	if err != nil {
		return points, err
	}

	return points, nil
}

func getNumComments(node *html.Node) (int, error) {
	num := 0
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return num, err
	}

	commentsStr := ""
	commentsQuery := htmlquery.Find(node, "/a]")
	for _, query := range commentsQuery {
		linkStr := htmlquery.InnerText(query)
		// No comments added to post
		if strings.Contains(linkStr, "discuss") {
			return 0, nil
		} else if strings.Contains(linkStr, "comment") {
			// Extract the number of comments
			commentsStr = reg.ReplaceAllLiteralString(linkStr, "")
		}
	}
	if commentsStr == "" {
		return num, errors.New(errorMsg)
	}

	num, err = strconv.Atoi(commentsStr)
	if err != nil {
		return num, err
	}

	return num, nil
}

func getTimePosted(node *html.Node) (time.Time, error) {
	var posted time.Time
	timeQuery := htmlquery.Find(node, "/span[contains(@class, 'age')]")
	if len(timeQuery) != 1 {
		return posted, errors.New(errorMsg)
	}
	timeStr := htmlquery.SelectAttr(timeQuery[0], "title")
	posted, err := time.Parse("2006-01-02T15:04:05", timeStr)
	if err != nil {
		return posted, err
	}

	return posted, nil
}
