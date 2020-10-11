package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"github.com/mmcdole/gofeed"

	"github.com/shinshin86/qiita-tag-feed-reader-cli/tag"
)

// MaxDisplayLen is maximum number of characters for content display
const MaxDisplayLen = 200

// FeedData store fetch feed data
type FeedData struct {
	Title       string
	FeedType    string
	FeedVersion string
	Items       []FeedItem
}

// FeedItem store fetch feed item
type FeedItem struct {
	Title       string
	Content     string
	Link        string
	Author      string
	PublishedAt string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func flagUsage() {
	usageText := `CLI reader of Qiita tag feed.

	Usage:
	qiita-tag-feed-reader-cli
	qiita-tag-feed-reader-cli <tag>
	qiita-tag-feed-reader-cli help`

	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func (f FeedItem) displayContent() string {
	if len(f.Content) < MaxDisplayLen {
		displayLen := len(f.Content)
		return f.Content[:displayLen]
	}

	return f.Content[:MaxDisplayLen]
}

func removeHTMLTag(html string) string {
	const pattern = `(<\/?[a-zA-A!-]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(html, -1)

	// Replace the long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})

	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			html = strings.ReplaceAll(html, group, "")
		}
	}
	return html
}

func removeNewlineTag(text string) string {
	return strings.NewReplacer(
		"\r\n", "",
		"\r", "",
		"\n", "",
	).Replace(text)
}

/*
	It runs on the CLI. So it display need like this.
	----------------------------------------------------
	<feed items>
	-----------------------
	<feed items>
	-----------------------
	<feed items>
	-----------------------
	<feed title>
	<feed type> <feed version>
	======================
*/
func contentDisplay(f FeedData) {
	sort.Slice(f.Items, func(i, j int) bool {
		return f.Items[i].PublishedAt < f.Items[j].PublishedAt
	})

	for _, item := range f.Items {
		fmt.Println(Bold(BrightCyan(item.Title)))
		fmt.Println("\t->", item.displayContent())
		fmt.Println("\t->", item.Link)
		fmt.Println("\t->", item.Author)
		fmt.Println("\t->", item.PublishedAt)
		fmt.Println("-----------------------")
	}

	fmt.Println(f.Title)
	fmt.Println(f.FeedType, f.FeedVersion)
	fmt.Println("======================")
}

func main() {
	flag.Usage = flagUsage

	var feedData FeedData
	var items []FeedItem
	var qtag string

	if len(os.Args) == 1 {
		qtag = tag.List[rand.Intn(100)]
	} else if len(os.Args) == 2 && os.Args[1] == "help" {
		flag.Usage()
		os.Exit(0)
	} else if len(os.Args) == 2 {
		qtag = os.Args[1]
	} else {
		fmt.Println("Invalid option")
		fmt.Println("==============")
		flag.Usage()
		os.Exit(1)
	}

	feedURL := "https://qiita.com/tags/" + qtag + "/feed.atom"
	feed, err := gofeed.NewParser().ParseURL(feedURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	feedData.Title = feed.Title
	feedData.FeedType = feed.FeedType
	feedData.FeedVersion = feed.FeedVersion

	for _, item := range feed.Items {
		if item == nil {
			break
		}

		items = append(items, FeedItem{
			Title:       item.Title,
			Content:     removeNewlineTag(removeHTMLTag(item.Content)),
			Link:        item.Link,
			Author:      item.Author.Name,
			PublishedAt: item.PublishedParsed.Format(time.RFC3339),
		})
	}

	feedData.Items = items
	contentDisplay(feedData)
}
