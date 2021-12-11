package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/c_anaerobic"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
)

func ListAnaerobicS(c *gin.Context) {
	response := app.NewResponse(c)

	svc := c_anaerobic.NewAnaerobicService(c, global.DBEngine, global.CacheEngine)
	resp, err := svc.ListAnaerobic()
	if err != nil {
		global.Logger.Errorf(c, "svc.NewAnaerobicService err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetAnaerobicS)
		return
	}

	response.ToResponseHtml("anaerobic.html", gin.H{
		"resp":      resp.Items,
		"sum":       resp.Sum,
		"tableName": "datatableRuns",
	})
}

//
func ListAnaerobicTimeZone(c *gin.Context) {
	response := app.NewResponse(c)
	typ := mtime.NewTimeTypeByStr(c.Param("typ"))
	svc := c_anaerobic.NewAnaerobicService(c, global.DBEngine, global.CacheEngine)
	resp, err := svc.ListAnaerobicZone(typ)
	if err != nil {
		global.Logger.Errorf(c, "svc.GerAnaerobicZoneFromDB err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetAnaerobicS)
		return
	}

	response.ToResponseHtml("anaerobic_month.html", gin.H{
		"resp":      resp.Items,
		"sum":       resp.Sum,
		"tableName": "datatableRuns",
	})
}
