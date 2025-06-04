package webmention

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/axxuy/webmention-sender/util"
	"golang.org/x/net/html"
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

func parsePage(page io.Reader) string {
	pageHtml, err := html.Parse(page)
	if err != nil {
		return ""
	}
	linkNodes, err := util.FindAllByTag(*pageHtml, []string{"a", "link"})
	if err != nil {
		return ""
	}
	for _, node := range linkNodes {
		var link string
		var hasRel bool
		for _, attr := range node.Attr {
			if attr.Key == "rel" && attr.Val == "webmention" {
				hasRel = true
			}
			if attr.Key == "href" {
				link = attr.Val
			}
		}
		if hasRel && link != "" {
			return link
		}
	}
	return ""
}
func GetWebmentionEndpoint(targetUrl *url.URL) (*Endpoint, error) {
	client := &http.Client{}
	//Check Header
	//headResp, err := client.Head(targetUrl.String())
	req, err := util.MakeRequest("HEAD", targetUrl.String())
	if err != nil {
		return nil, err
	}
	headResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer headResp.Body.Close()
	if headResp.StatusCode != http.StatusOK {
		return nil, nil
	}
	endpointUrl := parseLinkHeader(headResp.Header.Values("Link"))
	if endpointUrl != "" {
		url := util.ParseLink(endpointUrl)
		if url == nil {
			return nil, errors.New("Relative webmention endpoint")
		}
		return &Endpoint{client, url}, nil
	}

	//If there was nothing in the HEAD we'll need to GET the full page
	req, err = util.MakeRequest("GET", targetUrl.String())
	if err != nil {
		return nil, err
	}
	bodyResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer bodyResp.Body.Close()
	if bodyResp.StatusCode != http.StatusOK {
		return nil, nil
	}
	endpointUrl = parsePage(bodyResp.Body)
	if endpointUrl != "" {
		url := util.ParseLink(endpointUrl)
		if url == nil {
			return nil, errors.New("Relative webmention endpoint")
		}
		return &Endpoint{client, url}, nil
	}
	return nil, nil
}

func (e *Endpoint) SendWebmention(endpointUrl, targetUrl, sourceUrl *url.URL) error {
	if e == nil {
		return errors.New("Endpoint is nil")
	}
	if e.client == nil {
		return errors.New("Endpoint has no http Client")
	}
	if e.url == nil {
		return errors.New("Endpoint has no url")
	}
	body := url.Values{}
	body.Set("source", sourceUrl.String())
	body.Set("target", targetUrl.String())
	req, err := util.PostForm(e.url.String(), body)
	if err != nil {
		return err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	return nil
}
