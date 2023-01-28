package sdarot

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func (client *Client) login(username string, password string) error {
	initResponse, err := client.client.Get(SdarotURL)
	if err != nil {
		return errors.New("couldn't reach " + SdarotURL + ", the website might be down")
	}
	defer initResponse.Body.Close()

	endpoint, _ := url.JoinPath(SdarotURL, "login")

	params := url.Values{}
	params.Set("username", username)
	params.Set("password", password)
	params.Set("location", "%2F")
	params.Set("submit_login", "")

	response, err := client.client.PostForm(endpoint, params)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return fmt.Errorf("got invalid HTML: %w", err)
	}

	s := doc.Find(`.alert`)

	wrongCredsMsg := strings.TrimSpace(s.Text())
	if wrongCredsMsg != "" {
		return fmt.Errorf("%s", wrongCredsMsg)
	}

	return nil
}
