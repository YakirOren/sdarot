package sdarot

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context/ctxhttp"
)

func (client *Client) createSeason(ctx context.Context, seriesID int, season int) ([]VideoRequest, error) {
	episodesCount, err := client.fetchSeriesEpisodes(ctx, seriesID, season)
	if err != nil {
		return nil, err
	}

	var currentSeason []VideoRequest

	for j := 1; j < episodesCount+1; j++ {
		currentSeason = append(currentSeason, VideoRequest{
			SeriesID: seriesID,
			Season:   season,
			Episode:  j,
		})
	}

	return currentSeason, nil
}

func (client *Client) fetchSeriesEpisodes(ctx context.Context, id int, season int) (int, error) {
	watchURL := fmt.Sprintf("%s/ajax/watch?episodeList=%d&season=%d", SdarotURL, id, season)

	response, err := ctxhttp.Get(ctx, client.client, watchURL)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}

	defer response.Body.Close()

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read the response body: %w", err)
	}

	return strings.Count(string(all), "data-episode"), nil
}

func (client *Client) getWatchDoc(ctx context.Context, id int) (*goquery.Document, error) {
	response, err := ctxhttp.Get(ctx, client.client, SdarotURL+"/watch/"+strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("got invalid HTML: %w", err)
	}

	return doc, nil
}

func getSeasonsCount(doc *goquery.Document) int {
	const selector = "#season > li:nth-child(1)"
	seasonsCount := len(doc.Find(selector).NextAll().Nodes) + 1

	return seasonsCount
}

func extractNames(doc *goquery.Document) (string, string) {
	const selector = `.poster > div:nth-child(1) > h1:nth-child(1) > strong:nth-child(1)`
	titleElement := doc.Find(selector)

	data := strings.Split(titleElement.Text(), "/")
	hebrewName := strings.TrimSpace(data[0])
	englishName := strings.TrimSpace(data[1])

	return hebrewName, englishName
}
