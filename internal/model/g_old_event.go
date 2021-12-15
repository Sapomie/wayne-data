package model

type OldEvent struct {
	ID         uint `gorm:"primary_key"`
	Date       string
	TaskName   string
	Duration   float64
	Comment    string
	ParentTask string
	StartTime  int64
	EndTime    int64
}

func (e *OldEvent) TableName() string {
	return "g_old_event"
}

type OldEvents []*OldEvent
