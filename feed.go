package feed

import (
	"encoding/xml"
	"errors"
)

type atomFeed struct {
	Title     string `xml:"title"`
	Author    string `xml:"author"`
	NameSpace string `xml:"xmlns,attr"`
}
type atomEntry struct {
	Link    entryLink `xml:"link"` //This bit of inderection is needed to get the href. It seems go's xml parser is a bit fiddly with self closing tags
	Content string    `xml:"content"`
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
