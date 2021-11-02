package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListEvents(c *gin.Context) {
	param := service.EventListRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	svc := service.New(c)
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}

	casts, num, err := svc.GetEventList(&param, &pager)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetMovieList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEventListFail)
		return
	}

	response.ToResponseList(casts, num)

	//response.ToResponseHtml("a_casts.html", gin.H{
	//	"casts": casts,
	//})

}
