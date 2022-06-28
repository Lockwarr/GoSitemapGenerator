package scraper_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"
	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/scraper"

	"github.com/stretchr/testify/suite"
)

type scraperTestSuite struct {
	suite.Suite
	parallel   int
	maxDepth   int
	outputFile string
	scraper    scraper.ScraperService
}

func (s *scraperTestSuite) SetupTest() {
	s.parallel = 10
	s.maxDepth = 10
	s.outputFile = "testData/sitemap.xml"
	s.scraper = scraper.NewScraper(s.parallel, s.maxDepth, s.outputFile)
}

func (s *scraperTestSuite) AfterTest(suite string, testName string) {
}

func TestScraperTestSuite(t *testing.T) {
	suite.Run(t, &scraperTestSuite{})
}

func (s *scraperTestSuite) TestScrape_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	urlGenerated, _ := url.Parse("http://youtube.com")
	document, err := helpers.GetHtmlNode(ctx, urlGenerated)
	s.NoError(err)

	//Act
	gatheredUrls, err := s.scraper.Scrape(ctx, document)
	s.NoError(err)

	// Assert
	s.Equal(372, len(gatheredUrls))
}
