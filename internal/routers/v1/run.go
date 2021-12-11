package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/c_run"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
)

func ListRuns(c *gin.Context) {
	response := app.NewResponse(c)

	svc := c_run.NewRunService(c, global.DBEngine, global.CacheEngine)
	run, err := svc.ListRuns()
	if err != nil {
		global.Logger.Errorf(c, "svc.NewRunService err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetRuns)
		return
	}

	response.ToResponseHtml("run.html", gin.H{
		"resp":      run.Items,
		"avg":       run.Sum,
		"tableName": "datatableRuns",
	})
}

func ListRunTimeZone(c *gin.Context) {
	response := app.NewResponse(c)
	typ := mtime.NewTimeTypeByStr(c.Param("typ"))
	svc := c_run.NewRunService(c, global.DBEngine, global.CacheEngine)
	resp, err := svc.ListRunZone(typ)
	if err != nil {
		global.Logger.Errorf(c, "svc.NewRunService err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetRuns)
		return
	}

	response.ToResponseHtml("run_month.html", gin.H{
		"resp":      resp.Items,
		"avg":       resp.Sum,
		"tableName": "datatableRuns",
	})
}
