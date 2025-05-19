package util

import (
	"errors"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Find_all_by_tag(doc html.Node, tag string) ([]*html.Node, error) {
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
