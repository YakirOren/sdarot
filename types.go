package sdarot

import "net/url"

type Config struct {
	Username string
	Password string
}

type VideoRequest struct {
	SeriesID string
	Season   string
	Episode  string
}

type Video struct {
	ID       string
	URL      url.URL
	Metadata Metadata
}

type Metadata struct {
	Season      string
	Episode     string
	HebrewName  string
	EnglishName string
}
