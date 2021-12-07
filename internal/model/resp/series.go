package resp

type SeriesResp struct {
	Name          string
	Category      string
	Season        string
	Year          int
	EpisodeNumber int
	Duration      float64
	Rate          int
	Finish        string
	FirstTime     string
	LastTime      string
}

type SeriesSumResp struct {
	SeriesNumber int
	DurationAvg  float64
	RateAvg      int
	Finish       int
}
