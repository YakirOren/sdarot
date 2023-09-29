package sdarot

import (
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

const SdarotURL = "https://sdarot.tw"

type Client struct {
	client   *http.Client
	isMember bool
}

// New create a Client.
func New(config Config) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar:       jar,
		Transport: &refererTransport{},
	}

	sdarotClient := &Client{client: client, isMember: config.IsMember}

	err = sdarotClient.login(config.Username, config.Password)
	if err != nil {
		return nil, err
	}

	return sdarotClient, nil
}

type refererTransport struct{}

// RoundTrip adds a "Referer" header to every request sent from the Client.
func (t *refererTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Referer", SdarotURL)

	return http.DefaultTransport.RoundTrip(req)
}
