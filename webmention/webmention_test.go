package webmention

import "testing"

func Test_parseLinkHeader(t *testing.T) {
	//Empty header
	testEndpoint := "http://example.com/webmention"
	headers := []string{""}
	result := parseLinkHeader(headers)
	if result != "" {
		t.Error("Empty header should have empty result")
	}
	//Link header but no endpoint
	headers = []string{`<http://example.com/style.css>; rel=stylesheet`}
	result = parseLinkHeader(headers)
	if result != "" {
		t.Errorf("Expected %q, got %q\n", "", result)
	}

	//Webmention link but no rel
	headers = []string{`<http://example.com/webmention>; rel=stylesheet`}
	result = parseLinkHeader(headers)
	if result != "" {
		t.Errorf("Expected %q, got %q\n", "", result)
	}

	//Simple webmention link
	headers = []string{`<http://example.com/webmention>; rel=webmention`}
	result = parseLinkHeader(headers)
	if result != testEndpoint {
		t.Errorf("Expected %q, got %q\n", testEndpoint, result)
	}
	//Simple webmention link with quoted rel
	headers = []string{`<http://example.com/webmention>; rel="webmention"`}
	result = parseLinkHeader(headers)
	if result != testEndpoint {
		t.Errorf("Expected %q, got %q\n", testEndpoint, result)
	}
	// webmention link with query param
	headers = []string{`<http://example.com/webmention/?rel=webmention>`}
	result = parseLinkHeader(headers)
	if result != "" {
		t.Errorf("Expected %q, got %q\n", testEndpoint, result)
	}

	//Multiple links in header
	headers = []string{`<http://example.com/foo>, <http://example.com/webmention>; rel=webmention`}
	result = parseLinkHeader(headers)
	if result != testEndpoint {
		t.Errorf("Expected %q, got %q\n", testEndpoint, result)
	}

}
