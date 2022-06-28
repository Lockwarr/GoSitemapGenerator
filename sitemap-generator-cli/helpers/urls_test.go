package helpers_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/Lockwarr/GoSitemapGenerator/sitemap-generator-cli/helpers"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestValiateURL(t *testing.T) {
	// Arrange
	tests := []struct {
		name             string
		urlForValidation string
		expectedHost     string
		wantedErr        error
	}{
		{
			name:             "valid url",
			urlForValidation: "http://google.com",
			expectedHost:     "google.com",
			wantedErr:        nil,
		},
		{
			name:             "bad url",
			urlForValidation: "httpcom",
			expectedHost:     "google",
			wantedErr:        helpers.ErrInvalidURL,
		},
		{
			name:             "url without host",
			urlForValidation: "http://",
			expectedHost:     "google",
			wantedErr:        helpers.ErrURLNoSchemeOrHost,
		},
		{
			name:             "valid url different format",
			urlForValidation: "http://www.google.com/",
			expectedHost:     "www.google.com",
			wantedErr:        nil,
		},
		{
			name:             "valid url different format",
			urlForValidation: "http://www.google.bg/",
			expectedHost:     "www.google.bg",
			wantedErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			parsedUrl, err := helpers.ValidateURL(tt.urlForValidation)

			// Assert
			assert.Equal(t, tt.wantedErr, err)
			if tt.wantedErr == nil {
				assert.Equal(t, tt.expectedHost, parsedUrl.Host)
			}
		})
	}
}

func TestGetHtmlNode(t *testing.T) {
	// Arrange
	tests := []struct {
		name             string
		urlForValidation string
		expectedNodeType html.NodeType
		wantedErr        error
	}{
		{
			name:             "valid url",
			urlForValidation: "http://google.com",
			expectedNodeType: html.NodeType(0x2),
			wantedErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.ParseRequestURI(tt.urlForValidation)
			assert.NoError(t, err)

			// Act
			htmlNode, err := helpers.GetHtmlNode(context.Background(), parsedURL)
			// Assert
			assert.Equal(t, tt.wantedErr, err)
			assert.Equal(t, tt.expectedNodeType, htmlNode.Type)
		})
	}
}

func TestGatherURLS(t *testing.T) {
	// Arrange
	tests := []struct {
		name                   string
		urlForValidation       string
		expectedGatheredLength int
	}{
		{
			name:                   "valid url 1",
			urlForValidation:       "http://google.com",
			expectedGatheredLength: 7,
		},
		{
			name:                   "valid url 2",
			urlForValidation:       "http://youtube.com",
			expectedGatheredLength: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.ParseRequestURI(tt.urlForValidation)
			assert.NoError(t, err)
			document, err := helpers.GetHtmlNode(context.Background(), parsedURL)
			assert.NoError(t, err)

			// Act
			urls := helpers.GatherURLS(context.Background(), document, 10)

			// Assert
			assert.Equal(t, tt.expectedGatheredLength, len(urls))
		})
	}
}
