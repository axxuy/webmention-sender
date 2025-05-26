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
			attrPattern := regexp.MustCompile(`;.*rel\s*=\s*"?\s*webmention`)
			if attrPattern.MatchString(field) {
				link := util.CutSubString("<", ">", field)
				return link
			}
		}

	}
	return ""
}
func GetWebmentionEndpoint(targetUrl *url.URL) (*Endpoint, error) {
	client := &http.Client{}
	//Check Header
	resp, err := client.Head(targetUrl.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	endpointUrl := parseLinkHeader(resp.Header.Values("Link"))
	if endpointUrl == "" {
		return nil, nil
	}

	//If there was nothing in the HEAD we'll need to GET the full page
	return nil, nil

}

func (e *Endpoint) SendWebmention(endpointUrl, targetUrl, sourceUrl *url.URL) error {
	return nil
}
