package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/service/b_progress"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
)

func GetProgressNow(c *gin.Context) {
	response := app.NewResponse(c)

	typ := mtime.NewTimeTypeByStr(c.Param("typ"))
	zone := mtime.NewMTime(cons.DbNewest).TimeZone(typ)

	svc := b_progress.NewProgressService(c, global.DBEngine, global.CacheEngine)
	pro, err := svc.GetProgress(zone)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetProgress err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetProgressFail)
		return
	}
	response.ToResponseHtml("progress.html", gin.H{
		"progress": pro,
	})
}

func GetMobileProgressNow(c *gin.Context) {
	response := app.NewResponse(c)

	typ := mtime.NewTimeTypeByStr(c.Param("typ"))
	zone := mtime.NewMTime(cons.DbNewest).TimeZone(typ)

	svc := b_progress.NewProgressService(c, global.DBEngine, global.CacheEngine)
	pro, err := svc.GetProgress(zone)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetProgress err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetProgressFail)
		return
	}

	now := pro.GoalLefts[0]
	max := pro.GoalLefts[1]

	response.ToResponseHtml("progress_mobile.html", gin.H{
		"resp": pro,
		"now":  now,
		"max":  max,
	})
}
