package sdarot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	ErrServerOverLoad = errors.New("servers overload")
	ErrUnexpected     = errors.New("unexpected error")
)

const (
	overloadMsg = `כל שרתי הצפייה שלנו עמוסים ולא יכולים לטפל בפניות נוספות, נא לנסות שנית מאוחר יותר.<br />לצפייה בעומסי השרתים <a href="/status">לחצו כאן</a>.<br /><br /><b>משתמשים שתרמו לאתר יכולים לצפות בפרקים גם בשעות העומס!</b>`
	watchWait   = 30 * time.Second
	waitMsg     = "עליך להמתין 30 שניות, נא לבצע ריענון לעמוד זה"
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

	if !client.isMember {
		time.Sleep(watchWait)
	}

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
		return fmt.Errorf("faild to create download request: %w", err)
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
	params.Set("SID", strconv.Itoa(data.SeriesID))
	params.Set("ep", strconv.Itoa(data.Episode))
	params.Set("season", strconv.Itoa(data.Season))
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
	params.Set("serie", strconv.Itoa(data.SeriesID))
	params.Set("episode", strconv.Itoa(data.Episode))
	params.Set("season", strconv.Itoa(data.Season))
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
		if *output.Error == overloadMsg {
			return nil, ErrServerOverLoad
		}

		if *output.Error == waitMsg {
			// possible bug
			return nil, ErrServerOverLoad
		}

		return nil, fmt.Errorf("%w: %s", ErrUnexpected, *output.Error)
	}

	parsedURL, err := url.Parse(fmt.Sprintf("https:%s", output.Watch.URL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse returned url: %s", output.Watch.URL)
	}

	videoID, err := strconv.Atoi(output.VID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert: %w", err)
	}

	resp, err := client.client.Head(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("could not determine size of video")
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return &Video{
		ID:       videoID,
		URL:      *parsedURL,
		Metadata: data,
		Size:     resp.ContentLength,
	}, nil
}
