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
	MaximumParallel   int
	RecursionMaxDepth int
	OutputFilePath    string
	results           chan *url.URL
}

func NewScraper(parallel, maxDepth int, outputFile string) *Scraper {
	scraper := Scraper{
		OutputFilePath:    outputFile,
		MaximumParallel:   parallel,
		RecursionMaxDepth: maxDepth,
		results:           make(chan *url.URL),
	}
	return &scraper
}

// Scrape will scrape the given document and return all the urls found
func (s *Scraper) Scrape(ctx context.Context, document *html.Node) (urls []*url.URL, err error) {
	// we can modify GOMAXPROCS to run maxium parallel goroutines on different CPUs but
	// for Go versions > 1.15, the default value for GOMAXPROCS is number of available cores
	// so assuming the program is running at maximum speed
	var wg sync.WaitGroup
	wg.Add(s.MaximumParallel)
	urlsToProcess := make(chan *url.URL)

	for i := 0; i < s.MaximumParallel; i++ {
		go s.worker(ctx, urlsToProcess, &wg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		urlsFromCurrentDocument := helpers.GatherURLS(ctx, document, s.RecursionMaxDepth)
		for _, currentUrl := range urlsFromCurrentDocument {
			s.results <- currentUrl
			urlsToProcess <- currentUrl
		}
		close(urlsToProcess)
	}()

	go func() {
		wg.Wait()
		close(s.results)
	}()

	// wait for results
	var results []*url.URL
	for res := range s.results {
		results = append(results, res)
	}

	log.Println("Successfully colleted urls.")
	return results, nil
}

func (s *Scraper) worker(ctx context.Context, urlsToProcess chan *url.URL, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for currentUrl := range urlsToProcess {
			document, err := helpers.GetHtmlNode(ctx, currentUrl)
			if err != nil {
				// TODO: Could implement debug logger
				// log.Debug("failed to get html node for url: ", currentUrl, err)
				continue
			}
			urlsFromCurrentDocument := helpers.GatherURLS(ctx, document, s.RecursionMaxDepth)
			for _, currentUrl := range urlsFromCurrentDocument {
				s.results <- currentUrl
			}
		}
		return
	}
}
