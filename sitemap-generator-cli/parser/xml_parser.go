package parser

import (
	"encoding/xml"
	"net/url"
	"os"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"
)

// Sitemap is the root element of the XML file
type Sitemap struct {
	// XMLName is required. Encapsulates the file and references the current protocol standard.
	XMLName xml.Name `xml:"urlset"`
	//Urls is required. Parent tag for each URL entry. The remaining tags are children of this tag.
	Urls []SitemapURL `xml:"url"`
}

// SitemapURL is a single URL in a sitemap
type SitemapURL struct {
	// Loc required URL of the page. This URL must begin with the protocol (such as http) and end with a trailing slash, if your web server requires it. This value must be less than 2,048 characters.
	Loc string `xml:"loc"`
	// Lastmod is optional. TODO: implement
	LastMod string `xml:"lastmod"`
	// ChangeFreq is optional. TODO: implement
	ChangeFreq string `xml:"changefreq"`
	// Priority is optional. TODO: implement
	Priority string `xml:"priority"`
}

// UnmarshalXML unmarshals XML file into struct
func UnmarshalXML(filePath string, sitemap *Sitemap) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(sitemap)
	if err != nil {
		return err
	}

	return nil
}

// MarshalXML marshals struct into XML file
func MarshalXML(filePath string, sitemap *Sitemap) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("  ", "    ")
	err = encoder.Encode(sitemap)
	if err != nil {
		return err
	}

	return nil
}

// UrlsToSitemap - converts slice of URLs to Sitemap with urlset given url as first argument
func UrlsToSitemap(urlset *url.URL, urls []*url.URL) *Sitemap {
	var sitemapUrls []SitemapURL
	for _, url := range urls {
		// Sometimes we have urls such as (#learnmore, #advertise, //www.google.bg). We want to validate every url use only valid urls for the sitemap file.
		_, err := helpers.ValidateURL(url.String())
		if err != nil {
			continue
		}
		sitemapUrls = append(sitemapUrls, SitemapURL{
			Loc:        url.String(), // Loc is enough to satisfy the Sitemaps protocol. All other fields are optional.
			LastMod:    "",           // TODO:
			ChangeFreq: "",           // TODO:
			Priority:   "",           // TODO:
		})
	}

	return &Sitemap{XMLName: xml.Name{Space: urlset.String(), Local: urlset.Host}, Urls: sitemapUrls}
}
