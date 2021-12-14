package resp

type EventFieldResponse struct {
	Name      string
	FromNow   float64
	Longest   float64
	Duration  float64
	Times     int64
	FirstTime string
	LastTime  string

	LastTimeT int64
}
