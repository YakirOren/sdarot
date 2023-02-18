package sdarot

type Series struct {
	ID          int
	HebrewName  string
	EnglishName string
	Seasons     [][]VideoRequest
}

func (s *Series) BuildVideoRequest(season int, episode int) VideoRequest {
	return VideoRequest{
		SeriesID: s.ID,
		Season:   season,
		Episode:  episode,
	}
}

func (s *Series) GetSeasons() int {
	return len(s.Seasons)
}

func (s *Series) GetEpisodes(season int) []VideoRequest {
	return s.Seasons[season-1]
}
