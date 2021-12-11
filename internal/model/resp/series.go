package resp

type SeriesItem struct {
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

type SeriesSum struct {
	SeriesNumber int
	DurationAvg  float64
	RateAvg      int
	Finish       int
}

type Series struct {
	Item []*SeriesItem
	Sum  *SeriesSum
}
