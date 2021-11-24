package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/b_rawEvent"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ExportAllRawEvent(c *gin.Context) {
	response := app.NewResponse(c)

	svc := b_rawEvent.NewRawEventService(c, global.DBEngine, global.CacheEngine)
	err := svc.ExportAllRawEvent()
	if err != nil {
		global.Logger.Errorf(c, "svc.ExportAllRawEvent err: %v", err)
		response.ToErrorResponse(errcode.ErrorExportCsv)
		return
	}

	response.ToResponse("Success")

}
