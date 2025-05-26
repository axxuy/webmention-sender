package feed

import (
	"encoding/xml"
	"errors"
	"net/url"
	"strings"

	"github.com/axxuy/webmention-sender/util"
	"golang.org/x/net/html"
)

type atomFeed struct {
	Title     string      `xml:"title"`
	Author    string      `xml:"author"`
	NameSpace string      `xml:"xmlns,attr"`
	Entries   []atomEntry `xml:"entry"`
}
type atomEntry struct {
	Link    entryLink `xml:"link"` //This bit of inderection is needed to get the href. It seems go's xml parser is a bit fiddly with self closing tags
	Content string    `xml:"content"`
	Id      string    `xml:"id"`
	PubDate string    `xml:"published"`
}
type entryLink struct {
	Link string `xml:"href,attr"`
}

type Entry struct {
	Id    string
	Url   *url.URL
	Links []*url.URL
}

func parseLink(link string) *url.URL {
	url, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if !url.IsAbs() {
		return nil
	}
	return url
}

func convertEntry(entry atomEntry) (Entry, error) {
	postUrl, err := url.Parse(entry.Link.Link)
	if err != nil {
		return Entry{}, nil
	}
	contentHtml, err := html.Parse(strings.NewReader(entry.Content))
	if err != nil {
		return Entry{}, err
	}
	linkNodes, err := util.FindAllByTag(*contentHtml, []string{"a"})
	if err != nil {
		return Entry{}, err
	}
	links := make([]*url.URL, 0)
	for _, node := range linkNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				url := parseLink(attr.Val)
				if url == nil {
					continue
				}
				if url.Hostname() == postUrl.Hostname() {
					continue
				}
				links = append(links, url)
			}
		}
	}
	return Entry{entry.Id, postUrl, links}, nil
}

func ParseAtomFeed(data []byte) ([]Entry, error) {
	var feed atomFeed
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	if feed.NameSpace != "http://www.w3.org/2005/Atom" {
		return nil, errors.New("Invalid atom feed") //About the laziest kind of validation we can do, but sometimes the bar really is that low
	}
	numEntries := len(feed.Entries)
	entries := make([]Entry, 0, numEntries)
	for _, rawEntry := range feed.Entries {
		entry, err := convertEntry(rawEntry)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
