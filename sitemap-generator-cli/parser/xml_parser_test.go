package parser_test

import (
	"testing"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/parser"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalXML(t *testing.T) {
	// Arrange
	var sitemap parser.Sitemap
	expectedSitemapUrls := []parser.SitemapURL{{Loc: "http://www.example.com/", LastMod: "2005-01-01", ChangeFreq: "monthly", Priority: "0.8"}}
	filePath := "./testdata/sitemap.xml"

	// Act
	err := parser.UnmarshalXML(filePath, &sitemap)

	// Assert
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedSitemapUrls, sitemap.Urls)
}

func TestUnmarshal_WhenNonExistingFile_ThenFail(t *testing.T) {
	// Arrange
	var sitemap parser.Sitemap
	filePath := "12"

	// Act
	err := parser.UnmarshalXML(filePath, &sitemap)

	// Assert
	assert.Equal(t, "open 12: no such file or directory", err.Error())
}

func TestUnmarshalXML_WhenDecodingError_ThenFail(t *testing.T) {
	// Arrange
	var sitemap parser.Sitemap
	filePath := "./testdata/fakeSitemap.xml"

	// Act
	err := parser.UnmarshalXML(filePath, &sitemap)

	// Assert
	assert.Equal(t, "expected element type <urlset> but have <url>", err.Error())
}

func TestMarshalXML(t *testing.T) {
	// Arrange
	sitemap := &parser.Sitemap{
		Urls: []parser.SitemapURL{
			{Loc: "http://www.example.com/", LastMod: "2005-01-01", ChangeFreq: "monthly", Priority: "0.8"},
		},
	}
	outputFilePath := "./testdata/testResultSitemap.xml"

	// Act
	err := parser.MarshalXML(outputFilePath, sitemap)

	// Assert
	assert.Equal(t, nil, err)
}

func TestMarshalXML_WhenNonExistingFile_ThenFail(t *testing.T) {
	// Arrange
	sitemap := &parser.Sitemap{
		Urls: []parser.SitemapURL{
			{Loc: "http://www.example.com/", LastMod: "2005-01-01", ChangeFreq: "monthly", Priority: "0.8"},
		},
	}
	outputFilePath := ""

	// Act
	err := parser.MarshalXML(outputFilePath, sitemap)

	// Assert
	assert.Equal(t, "open : no such file or directory", err.Error())
}
