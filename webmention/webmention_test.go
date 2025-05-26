package webmention

import "testing"

func Test_parseLinkHeader(t *testing.T) {
	headers := []string{""}
	if parseLinkHeader(headers) != "" {
		t.Error("Empty header should have empty result")
	}
}
