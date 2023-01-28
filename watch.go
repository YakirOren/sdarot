package sdarot

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetVideo fetches the requested episode from Sdarot and returns a Video object.
// Note: there is a 30-second delay internally.
func (client *Client) GetVideo(data VideoRequest) (*Video, error) {
	return client.GetVideoWithContext(context.Background(), data)
}

// GetVideoWithContext fetches the requested episode from Sdarot and returns a Video object.
// Note: there is a 30-second delay internally.
func (client *Client) GetVideoWithContext(ctx context.Context, data VideoRequest) (*Video, error) {
	token, err := client.preWatch(ctx, data)
	if err != nil {
		return nil, err
	}

	const timeout = 30 * time.Second
	time.Sleep(timeout)

	return client.watch(ctx, data, token)
}

// Download the given video and write it into the given writer.
func (client *Client) Download(video *Video, writer io.Writer) error {
	return client.DownloadWithContext(context.Background(), video, writer)
}

// DownloadWithContext the given video and write it into the given writer.
func (client *Client) DownloadWithContext(ctx context.Context, video *Video, writer io.Writer) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, video.URL.String(), nil)
	if err != nil {
		return err
	}

	response, err := client.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed getting video: %w", err)
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return fmt.Errorf("faild to write response into the given writer: %w", err)
	}

	return nil
}

// preWatch returns a token that is required in order to fetch other resources.
func (client *Client) preWatch(ctx context.Context, data VideoRequest) (string, error) {
	endpoint := fmt.Sprintf(SdarotURL + "/ajax/watch")

	params := url.Values{}
	params.Set("SID", data.SeriesID)
	params.Set("ep", data.Episode)
	params.Set("season", data.Season)
	params.Set("preWatch", "true")

	response, err := ctxhttp.PostForm(ctx, client.client, endpoint, params)
	if err != nil {
		return "", fmt.Errorf("failed to get prewatch token: %w", err)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("reading token failed: %w", err)
	}

	return string(bytes), nil
}

// watch returns a Video for a given VideoRequest.
func (client *Client) watch(ctx context.Context, data VideoRequest, watchToken string) (*Video, error) {
	endpoint := fmt.Sprintf(SdarotURL + "/ajax/watch")

	params := url.Values{}
	params.Set("serie", data.SeriesID)
	params.Set("episode", data.Episode)
	params.Set("season", data.Season)
	params.Set("type", "episode")
	params.Set("watch", "false")
	params.Set("token", watchToken)

	response, err := ctxhttp.PostForm(ctx, client.client, endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("could not retch watch endpoint: %w", err)
	}

	defer response.Body.Close()

	var output struct {
		VID   string `json:"VID,omitempty"`
		Watch struct {
			URL string `json:"480"`
		} `json:"watch,omitempty"`
		Error *string `json:"error,omitempty"`
	}

	err = json.NewDecoder(response.Body).Decode(&output)
	if err != nil {
		return nil, fmt.Errorf("got invalid json: %w", err)
	}

	if output.Error != nil {
		return nil, fmt.Errorf("%s", *output.Error)
	}

	parsedURL, err := url.Parse(fmt.Sprintf("https:%s", output.Watch.URL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse returned url: %s", output.Watch.URL)
	}

	heb, eng, err := client.getSeriesName(ctx, data.SeriesID)
	if err != nil {
		return nil, fmt.Errorf("could not get the series name: %w", err)
	}

	return &Video{
		ID:  output.VID,
		URL: *parsedURL,
		Metadata: Metadata{
			Season:      data.Season,
			Episode:     data.Episode,
			HebrewName:  heb,
			EnglishName: eng,
		},
	}, nil
}

// getSeriesName returns the hebrew and english name for a given series.
func (client *Client) getSeriesName(ctx context.Context, seriesID string) (string, string, error) {
	response, err := ctxhttp.Get(ctx, client.client, SdarotURL+"/watch/"+seriesID)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute request: %w", err)
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("got invalid HTML: %w", err)
	}

	const selector = `.poster > div:nth-child(1) > h1:nth-child(1) > strong:nth-child(1)`
	titleElement := doc.Find(selector)

	data := strings.Split(titleElement.Text(), "/")
	hebrewName := strings.TrimSpace(data[0])
	englishName := strings.TrimSpace(data[1])

	return hebrewName, englishName, nil
}
