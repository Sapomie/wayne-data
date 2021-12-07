package resp

type BookResp struct {
	Name       string
	Category   string
	Author     string
	Year       int
	WordNumber float64
	Duration   float64
	Rate       int
	Finish     string
	FirstTime  string
	LastTime   string
}

type BookSumResp struct {
	BookNumber     int
	CategoryNumber int
	WordNumberAvg  float64
	DurationAvg    float64
	RateAvg        int
	Finish         int
}
