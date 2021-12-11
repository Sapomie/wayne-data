package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/c_series"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListSeriesS(c *gin.Context) {
	response := app.NewResponse(c)

	svc := c_series.NewSeriesService(c, global.DBEngine, global.CacheEngine)
	resp, err := svc.ListSeries()
	if err != nil {
		global.Logger.Errorf(c, "svc.NewSeriesService err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetSeries)
		return
	}

	response.ToResponseHtml("series.html", gin.H{
		"resp": resp.Item,
		"sum":  resp.Sum,
	})
}
