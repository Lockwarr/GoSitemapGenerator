package main

import (
	"flag"
	"log"
	"os"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"
	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/scraper"
)

func main() {
	log.Println("Starting sitemap generator cli tool...")

	if len(os.Args) <= 1 {
		log.Println("missing url argument, see readme")
		os.Exit(0)
	}

	site := os.Args[1]
	err := helpers.ValidateUrl(site)
	if err != nil {
		log.Println("invalid url", err)
		os.Exit(0)
	}
	parallel := flag.String("parallel", "10", "number of parallel workers to navigate through site")
	outputFile := flag.String("output-file", "sitemap.txt", "output file path")
	maxDepth := flag.String("max-depth", "10", "max depth of url naviagtion recursion")
	flag.Parse()
	scraperInstance := scraper.NewScraper(parallel, outputFile, maxDepth)
	scraperInstance.Scrape(site)

}
