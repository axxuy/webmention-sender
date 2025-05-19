package feed

import "encoding/xml"

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
	err := xml.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return nil, nil

}
