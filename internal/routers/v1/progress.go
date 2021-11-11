package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/service/progress"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
	"time"
)

func GetProgressTenNow(c *gin.Context) {
	response := app.NewResponse(c)
	zone := mtime.NewMTime(cons.Newest.Add(-1 * time.Hour)).TimeZone(mtime.TypeTen)
	startTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)

	svc := progress.NewEssentialService(c)
	pro, err := svc.GetProgress(zone, startTime)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetProgress err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetProgressFail)
		return
	}
	response.ToResponseHtml("progress.html", gin.H{
		"progress": pro,
	})
}

func GetProgressMonthNow(c *gin.Context) {
	response := app.NewResponse(c)

	typ := mtime.NewTimeTypeByStr(c.Param("typ"))

	zone := mtime.NewMTime(cons.Newest.Add(-1 * time.Hour)).TimeZone(typ)
	startTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)

	svc := progress.NewEssentialService(c)
	pro, err := svc.GetProgress(zone, startTime)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetProgress err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetProgressFail)
		return
	}
	response.ToResponseHtml("progress.html", gin.H{
		"progress": pro,
	})
}
