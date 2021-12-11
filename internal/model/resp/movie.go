package resp

type MovieItem struct {
	Date   string
	Name   string
	NameEn string
	Rate   int
	Year   int
	Place  string
}

type MovieSum struct {
	MovieNumber  int
	Rate         int
	CinemaNumber int
}

type Movie struct {
	Items []*MovieItem
	Sum   *MovieSum
}
