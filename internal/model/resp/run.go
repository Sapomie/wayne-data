package resp

type RunResp struct {
	Id          int
	Date        string
	Distance    float64
	TimeCost    float64
	Pace        string
	Speed       float64
	Rate        int
	Temperature int
	Altitude    int
}

type RunSum struct {
	Id                 int
	Date               string
	Times              int
	Distance           float64
	DistanceAverage    float64
	Pace               string
	Speed              float64
	RateAverage        int
	TemperatureAverage int
	AltitudeAverage    int
}
