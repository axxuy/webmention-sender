package feed

import (
	"encoding/xml"
	"errors"
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
	Url   string
	Links []string
}

func convertEntry(entry atomEntry) (Entry, error) {
	result := Entry{}
	result.Id = entry.Id
	result.Url = entry.Link.Link
	contentHtml, err := html.Parse(strings.NewReader(entry.Content))
	if err != nil {
		return Entry{}, err
	}
	linkNodes, err := util.FindAllByTag(*contentHtml, "a")
	if err != nil {
		return Entry{}, err
	}
	links := make([]string, 0)
	for _, node := range linkNodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	result.Links = links
	return result, nil
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
	return nil, nil
}
