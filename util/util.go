package util

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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

func FindAllByTag(doc html.Node, tag string) ([]*html.Node, error) {
	tagAtom := atom.Lookup([]byte(tag))
	if tagAtom == atom.Atom(0) {
		return nil, errors.New("No such tag")
	}
	result := make([]*html.Node, 0)
	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == tagAtom {
			result = append(result, node)
		}
	}
	return result, nil
}
