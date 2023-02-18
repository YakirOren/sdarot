package sdarot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

type SearchResult struct {
	SeriesID    string
	HebrewName  string
	EnglishName string
}

// Search returns an array of SearchResult for a given term.
func (client *Client) Search(term string) ([]*Series, error) {
	return client.SearchWithContext(context.Background(), term)
}

func (client *Client) GetSeriesByID(id int) (*Series, error) {
	return client.GetSeriesWithContext(context.Background(), id)
}

func (client *Client) GetSeriesWithContext(ctx context.Context, seriesID int) (*Series, error) {
	doc, err := client.getWatchDoc(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	seasonsCount := getSeasonsCount(doc)

	var Seasons [][]VideoRequest

	for i := 1; i < seasonsCount+1; i++ {
		newSeason, err := client.createSeason(ctx, seriesID, i)
		if err != nil {
			return nil, err
		}

		Seasons = append(Seasons, newSeason)
	}

	hebrewName, englishName := extractNames(doc)

	return &Series{
		ID:          seriesID,
		HebrewName:  hebrewName,
		EnglishName: englishName,
		Seasons:     Seasons,
	}, nil
}

// SearchWithContext returns an array of SearchResult for a given term.
func (client *Client) SearchWithContext(ctx context.Context, term string) ([]*Series, error) {
	endpoint := fmt.Sprintf(SdarotURL+"/ajax/index?search=%s", url.QueryEscape(term))

	response, err := ctxhttp.Get(ctx, client.client, endpoint)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
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

	searchResults := make([]*Series, len(rawSearchResults))

	for _, result := range rawSearchResults {
		names := strings.Split(result.Name, " / ")

		seriesID, err := strconv.Atoi(result.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert: %w", err)
		}

		searchResults = append(searchResults, &Series{
			seriesID,
			strings.TrimSpace(names[0]),
			strings.TrimSpace(names[1]),
			nil,
		})
	}

	return searchResults, nil
}
