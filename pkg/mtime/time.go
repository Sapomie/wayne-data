package mtime

import (
	"fmt"
	"time"
)

const (
	TimeTemplate1  = "2006/1/2 15:04:05"
	TimeTemplate2  = "2006/1/02 15:04:05"
	TimeTemplate3  = "2006/01/02 15:04:05"
	TimeTemplate4  = "2006-01-02"
	TimeTemplate5  = "2006-01-02 15:04"
	TimeTemplate6  = "01-02 15:04"
	TimeTemplate7  = "2006-01"
	TimeTemplate8  = "01-02"
	TimeTemplate9  = "15:04:05"
	TimeTemplate10 = "2006/01/02"
	DefaultTime    = "1990-10-27 20:00:00"
)

type TimeType int

const (
	TypeDay TimeType = 1 + iota
	TypeTen
	TypeMonth
	TypeQuarter
	TypeHalf
	TypeYear
)

func NewTimeTypeByStr(str string) TimeType {
	switch str {
	case "day":
		return TypeDay
	case "ten":
		return TypeTen
	case "month":
		return TypeMonth
	case "quarter":
		return TypeQuarter
	case "half":
		return TypeHalf
	case "year":
		return TypeYear
	default:
		return TypeTen
	}
}

type TimeZone struct {
	Typ  TimeType
	Year int
	Num  int
}

type TimeZoneInterface interface {
	BeginAndEnd() (begin, end time.Time)
	DateString() (date string)
	ZoneTyp() TimeType
	Days() (days float64)
	Number() (num int)
}

func NewTimeZone(typ TimeType, year, num int) *TimeZone {
	return &TimeZone{
		Typ:  typ,
		Year: year,
		Num:  num,
	}
}

func (tz *TimeZone) BeginAndEnd() (begin, end time.Time) {
	switch tz.Typ {
	case TypeDay:
		begin = time.Date(tz.Year, 1, tz.Num, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year, 1, tz.Num+1, 0, 0, 0, 0, time.Local)
	case TypeTen:
		begin = time.Date(tz.Year, 1, 1+(tz.Num-1)*10, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year, 1, 1+tz.Num*10, 0, 0, 0, 0, time.Local)
	case TypeMonth:
		month := time.Month(tz.Num)
		begin = time.Date(tz.Year, month, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year, month+1, 1, 0, 0, 0, 0, time.Local)
	case TypeQuarter:
		startMonth := time.Month((tz.Num-1)*3 + 1)
		begin = time.Date(tz.Year, startMonth, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year, startMonth+3, 1, 0, 0, 0, 0, time.Local)
	case TypeHalf:
		begin = time.Date(tz.Year, time.Month(1+(tz.Num-1)*6), 1, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year, time.Month(1+tz.Num*6), 1, 0, 0, 0, 0, time.Local)
	case TypeYear:
		begin = time.Date(tz.Year, 1, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(tz.Year+1, 1, 1, 0, 0, 0, 0, time.Local)
	}
	return
}

func (tz *TimeZone) DateString() (date string) {
	start, _ := tz.BeginAndEnd()
	mt := NewMTime(start)

	switch tz.Typ {
	case TypeDay:
		date = start.Format(TimeTemplate8) + "/" + mt.WeekDayShort()
	case TypeTen:
		date = "Ten" + fmt.Sprint(mt.Ten())
	case TypeMonth:
		date = start.Month().String()
	case TypeQuarter:
		date = "Quarter" + fmt.Sprint(mt.Quarter())
	case TypeHalf:
		date = "Half" + fmt.Sprint(mt.Half())
	case TypeYear:
		date = fmt.Sprint(start.Year())
	}

	return
}

func (tz *TimeZone) ZoneTyp() TimeType {
	return tz.Typ
}

func (tz *TimeZone) Days() (days float64) {
	start, end := tz.BeginAndEnd()
	return end.Sub(start).Hours() / 24
}

func (tz *TimeZone) Number() (num int) {
	return tz.Num
}

type MTime struct {
	Time time.Time
}

func NewMTime(t time.Time) *MTime {
	return &MTime{t}
}

func (mt *MTime) Day() (day int) {
	return mt.Time.YearDay()
}

func (mt *MTime) Ten() (ten int) {
	beginningOfTheYear, _ := NewTimeZone(TypeYear, mt.Time.Year(), 1).BeginAndEnd()
	days := (mt.Time.Sub(beginningOfTheYear).Hours()) / 24
	ten = int(days)/10 + 1
	return
}

func (mt *MTime) Month() (month int) {
	return int(mt.Time.Month())
}

func (mt *MTime) Quarter() (quarter int) {
	month := int(mt.Time.Month())
	quarter = month/3 + 1
	if month%3 == 0 {
		quarter = month / 3
	}
	return
}

func (mt *MTime) Half() int {
	month := int(mt.Time.Month())
	if month < 7 {
		return 1
	} else {
		return 2
	}
}

func (mt *MTime) TimeZone(typ TimeType) *TimeZone {
	var num int
	switch typ {
	case TypeDay:
		num = mt.Day()
	case TypeTen:
		num = mt.Ten()
	case TypeMonth:
		num = mt.Month()
	case TypeQuarter:
		num = mt.Quarter()
	case TypeHalf:
		num = mt.Half()
	case TypeYear:
		num = 1
	}

	return &TimeZone{
		Typ:  typ,
		Year: mt.Time.Year(),
		Num:  num,
	}

}

func (mt *MTime) WeekDayShort() string {
	return mt.Time.Weekday().String()[:3]
}

//保证start time 是一天的开始，end time 是一天的结尾
func TrimTime(startTime, endTime time.Time) (start, end time.Time) {

	start = trim(startTime)
	if start != startTime {
		start = start.Add(time.Hour * 24)
	}
	end = trim(endTime)

	//dur = end.Sub(start).Hours() / 24
	//if dur != float64(int(dur)) {
	//
	//}
	return
}

func trim(t time.Time) time.Time {
	return t.Add(time.Hour * 8).Truncate(time.Hour * 24).Add(-time.Hour * 8)
}
