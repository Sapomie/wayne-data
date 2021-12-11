package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	b_event2 "github.com/Sapomie/wayne-data/internal/service/b_event_field"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListEventField(c *gin.Context) {
	response := app.NewResponse(c)
	typ := model.NewFieldTypeByStr(c.Param("typ"))

	svc := b_event2.NewEvtFieldService(c)
	stuffs, err := svc.GetFieldList(typ)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetFieldList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEventField)
		return
	}

	response.ToResponseHtml("field.html", gin.H{
		"resp":      stuffs,
		"tableName": "datatableField",
	})

}
