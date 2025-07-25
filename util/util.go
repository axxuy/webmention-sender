package util

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const USER_AGENT = "Webmention Sender/1.0.1 (https://github.com/axxuy/webmention-sender)"

func MakeRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	//Set custom headers
	req.Header.Set("User-Agent", USER_AGENT)
	return req, nil
}
func PostForm(url string, body url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// return a substring of s between before and after. If either separator does not appear, return ""
func CutSubString(before, after, s string) string {
	_, suffix, found := strings.Cut(s, before)
	if !found {
		return ""
	}
	substr, _, found := strings.Cut(suffix, after)
	if !found {
		return ""
	}
	return substr

}

func FindAllByTag(doc html.Node, tags []string) ([]*html.Node, error) {
	atomTags := make([]atom.Atom, len(tags))
	for i, tag := range tags {
		tagAtom := atom.Lookup([]byte(tag))
		if tagAtom == atom.Atom(0) {
			return nil, errors.New("No such tag")
		}
		atomTags[i] = tagAtom
	}
	result := make([]*html.Node, 0)
	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && slices.Contains(atomTags, node.DataAtom) {
			result = append(result, node)
		}
	}
	return result, nil
}

func ParseLink(link string) *url.URL {
	if link == "" {
		return nil
	}
	url, err := url.Parse(link)
	if err != nil {
		return nil
	}
	if !url.IsAbs() {
		return nil
	}
	return url
}
