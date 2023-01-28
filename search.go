package sdarot

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"net/url"
	"strings"
)

type SearchResult struct {
	SeriesID    string
	HebrewName  string
	EnglishName string
}

// Search returns an array of SearchResult for a given term.
func (client *Client) Search(term string) ([]SearchResult, error) {
	return client.SearchWithContext(context.Background(), term)
}

// SearchWithContext returns an array of SearchResult for a given term.
func (client *Client) SearchWithContext(ctx context.Context, term string) ([]SearchResult, error) {
	endpoint := fmt.Sprintf(SdarotURL+"/ajax/index?search=%s", url.QueryEscape(term))

	response, err := ctxhttp.Get(ctx, client.client, endpoint)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var rawSearchResults []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err = json.NewDecoder(response.Body).Decode(&rawSearchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search result: %w", err)
	}

	var searchResults []SearchResult

	for _, result := range rawSearchResults {
		names := strings.Split(result.Name, " / ")

		searchResults = append(searchResults, SearchResult{
			SeriesID:    result.ID,
			HebrewName:  strings.TrimSpace(names[0]),
			EnglishName: strings.TrimSpace(names[1]),
		})
	}

	return searchResults, nil
}
