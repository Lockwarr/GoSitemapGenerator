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
func ValidateURL(urlForValidation string) error {
	parsedURL, err := url.ParseRequestURI(urlForValidation)
	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ErrURLNoSchemeOrHost
	}
	return nil
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
