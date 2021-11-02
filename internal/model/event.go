package model

type Event struct {
	//原有属性
	ID           int64   `gorm:"primary_key" json:"id"`
	Date         string  `gorm:"not null" json:"date"`
	Duration     float64 `gorm:"not null" json:"duration"`
	TaskNameId   int     `gorm:"not null;default:-2" json:"task_name_id"`
	ParentTaskId int     `gorm:"not null;default:-2" json:"parent_task_id"`
	Comment      string  `json:"comment"`
	StartTime    int64   `gorm:"not null" json:"start_time"`
	EndTime      int64   `gorm:"not null" json:"end_time"`
	//自定义属性：通过comment增加

	StuffId   string `json:"stuff_id"`
	TagId     string `json:"tag_id"`
	ProjectId int    `json:"project_id"`
	Remark    string `json:"remark"` //comment 除去自定义属性的部分

	*Model
}
