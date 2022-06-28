package parser

import (
	"encoding/xml"
	"os"
)

// Sitemap is the root element of the XML file
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	Urls    []SitemapURL `xml:"url"`
}

// SitemapURL is a single URL in a sitemap
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
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
