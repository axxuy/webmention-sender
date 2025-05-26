package webmention

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/axxuy/webmention-sender/util"
)

type Endpoint struct {
	client *http.Client
	url    *url.URL
}

// Check the http Link header for a url with rel="webmention"
// https://httpwg.org/specs/rfc8288.html#header
// Returns "" if none present
func parseLinkHeader(header []string) string {
	for _, headerVal := range header {
		if !strings.Contains(headerVal, "webmention") {
			continue
		}
		headerFields := strings.Split(headerVal, ",")
		for _, field := range headerFields {
			fieldParts := strings.Split(field, ";")
			for _, part := range fieldParts {
				attrPattern := regexp.MustCompile(`rel\s*=\s*"?\s*webmention`)
				if attrPattern.MatchString(part) {
					link := util.CutSubString("<", ">", part)
					return link
				}
			}
		}

	}
	return ""
}
func GetWebmentionEndpoint(targetUrl *url.URL) (*Endpoint, error) {
}

func (e *Endpoint) SendWebmention(endpointUrl, targetUrl, sourceUrl *url.URL) error {
}
