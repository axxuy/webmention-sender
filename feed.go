package feed

type Feed struct {
	Title     string `xml:"title"`
	Author    string `xml:"author"`
	NameSpace string `xml:"xmlns,attr"`
}
type Entry struct {
	Link    EntryLink `xml:"link"` //This bit of inderection is needed to get the href. It seems go's xml parser is a bit fiddly with self closing tags
	Content string    `xml:"content"`
}
type EntryLink struct {
	Link string `xml:"href,attr"`
}
