package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"
	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/parser"
	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/scraper"
)

func main() {
	log.Println("Starting sitemap generator cli tool...")
	ctx := context.Background()

	parallel := flag.Int("parallel", 10, "number of parallel workers to navigate through site")
	outputFile := flag.String("output-file", "../sitemap.xml", "output file path")
	maxDepth := flag.Int("max-depth", 10, "max depth of url naviagtion recursion")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Println("invalid number of arguments, see readme")
		os.Exit(0)
	}
	site := flag.Args()[0]

	url, err := helpers.ValidateURL(site)
	if err != nil {
		log.Println("invalid url", err)
		os.Exit(0)
	}

	scraperInstance := scraper.NewScraper(*parallel, *maxDepth, *outputFile)

	node, err := helpers.GetHtmlNode(ctx, url)
	if err != nil {
		log.Println("invalid url", err)
	}

	gatheredUrls, err := scraperInstance.Scrape(ctx, node)
	if err != nil {
		log.Println("error crawling through site", err)
	}

	err = parser.MarshalXML(*outputFile, parser.UrlsToSitemap(url, gatheredUrls))
	if err != nil {
		log.Println("error marshalling to xml", err)
	}

	log.Println("Successfully generated sitemap.xml, stopping program...")
}
