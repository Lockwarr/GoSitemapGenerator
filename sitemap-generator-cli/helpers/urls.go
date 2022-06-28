package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var (
	ErrInvalidURL        = errors.New("bad url entered as an argument")
	ErrURLNoSchemeOrHost = errors.New("url has no scheme or host")
)

// ValidateURL validates the url
func ValidateURL(urlForValidation string) (*url.URL, error) {
	parsedURL, err := url.ParseRequestURI(urlForValidation)
	if err != nil {
		return nil, ErrInvalidURL
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, ErrURLNoSchemeOrHost
	}
	return parsedURL, nil
}

// GetHtmlNode returns the html node of the given url
func GetHtmlNode(ctx context.Context, url *url.URL) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing the request %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing the response body %w", err)
	}

	return document, nil
}

func GatherURLS(ctx context.Context, n *html.Node, recursionMaxDepth int) []*url.URL {
	gatheredUrls := []*url.URL{}
	var f func(*html.Node, int) []*url.URL

	f = func(n *html.Node, depth int) []*url.URL {
		// There can only be one single <base> element in a document, and it must be inside the <head> element.
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "base") { // only get links from a and base tags
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "" {
					hrefURL, err := url.Parse(attr.Val)
					if err != nil {
						// TODO: Could implement debug logger
						// log.Debug("malformed href value: ", attr.Val, err)
						break // nested links are forbidden => assuming there is only one href per node, we can break from the loop
					}
					gatheredUrls = append(gatheredUrls, hrefURL)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if depth >= recursionMaxDepth {
				// max recursion depth reached, returning gathered links
				// could be modified with custom logger and log only on Debug level
				// log.Debug("max recursion depth reached, returning gathered links")
				return gatheredUrls
			}
			f(c, depth+1)
		}

		return gatheredUrls
	}
	// Starts first iteration
	f(n, 1)

	return gatheredUrls
}
