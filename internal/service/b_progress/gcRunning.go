package b_progress

import (
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/service/b_essential"
)

type GcRunning struct {
	RunningDistance float64
	GcAccumulation  int
	GcUsed          int
	GcLeft          int
}

func gcRunningInfo(es *b_essential.Essential) (info *GcRunning) {
	//get distance

	runningDistance := es.TaskInfo[cons.Running].Done * 5

	info = &GcRunning{
		RunningDistance: runningDistance,
		GcAccumulation:  int(runningDistance / 10),
		GcUsed:          int(es.StuffInfo[cons.StuGc].Done),
		GcLeft:          int(runningDistance/10) - int(es.StuffInfo[cons.StuGc].Done),
	}

	return
}
