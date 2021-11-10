package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/progress"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
	"time"
)

func GetProgress(c *gin.Context) {
	response := app.NewResponse(c)
	//param := resp.EventListRequest{}
	//valid, errs := app.BindAndValid(c, &param)
	//if !valid {
	//	global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
	//	response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
	//	return
	//}
	startTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
	svc := progress.NewEssentialService(c)
	pro, err := svc.GetProgress(mtime.TypeMonth, 2021, 11, startTime)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEventList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEventListFail)
		return
	}

	//response.ToResponse(casts)

	response.ToResponseHtml("summary.html", gin.H{
		"resp": pro,
	})

}
