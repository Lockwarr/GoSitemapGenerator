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
		wantedErr        error
	}{
		{
			name:             "valid url",
			urlForValidation: "http://google.com",
			wantedErr:        nil,
		},
		{
			name:             "bad url",
			urlForValidation: "httpcom",
			wantedErr:        helpers.ErrInvalidURL,
		},
		{
			name:             "url without host",
			urlForValidation: "http://",
			wantedErr:        helpers.ErrURLNoSchemeOrHost,
		},
		{
			name:             "valid url different format",
			urlForValidation: "http://www.google.com/",
			wantedErr:        nil,
		},
		{
			name:             "valid url different format",
			urlForValidation: "http://www.google.bg/",
			wantedErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := helpers.ValidateURL(tt.urlForValidation)

			// Assert
			assert.Equal(t, tt.wantedErr, err)
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
