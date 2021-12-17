package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/resp"
	"github.com/Sapomie/wayne-data/internal/service/g_old_event"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListOldEvents(c *gin.Context) {
	response := app.NewResponse(c)
	param := resp.OldEventListRequest{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := g_old_event.NewOldEventService(c)
	limit, offset := app.GetLimitOffset(c)
	events, num, err := svc.GetOldEventList(&param, limit, offset)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetOldEventList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEventListFail)
		return
	}

	response.ToResponseHtml("event_old.html", gin.H{
		"events": events,
		"num":    num,
	})

}
