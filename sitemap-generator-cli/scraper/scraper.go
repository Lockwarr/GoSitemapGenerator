package scraper

import (
	"context"
	"log"
	"net/url"
	"sync"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"

	"golang.org/x/net/html"
)

type ScraperService interface {
	Scrape(ctx context.Context, document *html.Node) (urls []*url.URL, err error)
}

type Scraper struct {
	limiter           helpers.ConcurrencyLimiter
	MaximumParallel   int
	RecursionMaxDepth int
	OutputFilePath    string
	// wg                *sync.WaitGroup
}

func NewScraper(parallel, maxDepth int, outputFile string) *Scraper {
	scraper := Scraper{
		OutputFilePath:    outputFile,
		MaximumParallel:   parallel,
		RecursionMaxDepth: maxDepth,
		limiter:           *helpers.NewConcurrencyLimiter(parallel, &sync.WaitGroup{}),
	}
	return &scraper
}

// Scrape extracts urls from a valid html document and starts worker for each found url
func (s *Scraper) Scrape(ctx context.Context, document *html.Node) (urls []*url.URL, err error) {
	urlsChan := make(chan *url.URL)
	var f func(*html.Node, int)

	f = func(n *html.Node, depth int) {
		if depth > s.RecursionMaxDepth {
			// max recursion depth reached, returning gathered links
			// could add Debug log here
			return
		}

		if s.limiter.GetNumInProgress() >= int32(s.MaximumParallel) {
			// max parallel workers reached, returning gathered links
			// could add Debug log here
			return
		}

		// we start nested recursion here, its' depth should be independent from the depth of the current recursion
		urlsFromCurrentDocument := s.gatherURLS(ctx, document, 1)
		for _, currentUrl := range urlsFromCurrentDocument {

			// start parallel worker for each url found
			s.limiter.Execute(func() {
				document, err := helpers.GetHtmlNode(ctx, currentUrl)
				if err != nil {
					log.Println("failed to get html node for url: ", currentUrl, err)
					return
				}

				if s.limiter.GetNumInProgress() >= int32(s.MaximumParallel) {
					// max parallel workers reached, returning gathered links
					// could add Debug log here
					return
				}
				urlsChan <- currentUrl
				f(document, depth+1)
			})
		}
	}
	f(document, 1)

	go func() {
		s.limiter.Wait()
		close(urlsChan)
	}()

	// wait for results
	var results []*url.URL
	for res := range urlsChan {
		results = append(results, res)
	}

	log.Println("Successfully colleted urls.")
	return results, nil
}

func (s *Scraper) gatherURLS(ctx context.Context, n *html.Node, recursionDepth int) []*url.URL {
	gatheredUrls := []*url.URL{}
	var f func(*html.Node, int) []*url.URL

	f = func(n *html.Node, depth int) []*url.URL {
		// the task says `take in account <base> element if declared`, but it's not clear
		// however I assume that if <base> element is declared, we should try and find more links there
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "base") { // only get links from a and base tags
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "" {
					hrefURL, err := url.Parse(attr.Val)
					if err != nil {
						log.Println("malformed href value: ", attr.Val, err)
						break // nested links are forbidden => assuming there is only one href per node, we can break from the loop
					}
					gatheredUrls = append(gatheredUrls, hrefURL)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if depth >= s.RecursionMaxDepth {
				// max recursion depth reached, returning gathered links
				// could be modified with custom logger and log only on Debug level
				// log.Debug("max recursion depth reached, returning gathered links")
				return gatheredUrls
			}
			f(c, depth+1)
		}

		return gatheredUrls
	}
	f(n, recursionDepth)

	return gatheredUrls
}
