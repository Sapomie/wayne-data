package resp

type AnaerobicItem struct {
	Date     string
	Name     string
	Group    int
	Times    int
	Addition float64
}

type AnaerobicSum struct {
	Id                    int
	Date                  string
	Protein1              int
	TotalGroup            int
	SitUpGroups           int
	SitUpPerGroup         int
	PushUpGroups          int
	PushUpPerGroup        int
	DumbbellPressGroups   int
	DumbbellPressPerGroup int
	DumbbellPressMass     float64
}

type Anaerobic struct {
	Items []*AnaerobicItem
	Sum   *AnaerobicSum
}

type AnaerobicZone struct {
	Items []*AnaerobicSum
	Sum   *AnaerobicSum
}
