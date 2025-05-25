package webmention

import (
	"net/http"
	"net/url"
)

type Endpoint struct {
	client *http.Client
	url    *url.URL
}

func GetWebmentionEndpoint(targetUrl *url.URL) (*Endpoint, error) {
}

func (e *Endpoint) SendWebmention(endpointUrl, targetUrl, sourceUrl *url.URL) error {
}
