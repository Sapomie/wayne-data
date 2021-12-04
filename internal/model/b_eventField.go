package model

const (
	TypeParent = iota + 1
	TypeTask
	TypeStuff
	TypeTag
	TypeProject
)

func NewFieldTypeByStr(str string) int {
	switch str {
	case "parent":
		return TypeParent
	case "task":
		return TypeTask
	case "stuff":
		return TypeStuff
	case "tag":
		return TypeTag
	case "project":
		return TypeProject
	default:
		return TypeStuff
	}
}

type EventField interface {
	FieldName() string
	FieldTotalDuration() float64
	FieldEventNum() int64
	FieldFirstTimeAndLastTime() (int64, int64)
	FieldLongest() int64
}

type EventFields interface {
	ToEventFields() []EventField
}
