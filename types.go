package sdarot

import "net/url"

type Config struct {
	Username string
	Password string
	IsMember bool
}

type VideoRequest struct {
	SeriesID int
	Season   int
	Episode  int
}

type Video struct {
	ID       int
	URL      url.URL
	Metadata VideoRequest
	Size     int64
}
