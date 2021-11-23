package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/a_procession"
	"github.com/Sapomie/wayne-data/internal/service/b_rawEvent"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	response := app.NewResponse(c)
	response.ToResponseHtml("upload.html", gin.H{})
}

func UploadPost(c *gin.Context) {
	response := app.NewResponse(c)

	//获取文件
	file, err := c.FormFile("myFile")
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadPost err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFail)
		return
	}
	//存储CSV文件
	err = c.SaveUploadedFile(file, global.AppSetting.CsvSavePath+file.Filename)
	if err != nil {
		global.Logger.Errorf(c, "svc.SaveUploadedFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorSaveUploadingFile)
		return
	}
	//处理文件
	_, importDataInfos, err := b_rawEvent.NewRawEventService(c, global.DBEngine, global.CacheEngine).ImportCsvData()
	if err != nil {
		global.Logger.Errorf(c, "svc.NewRawEventService err: %v", err)
		response.ToErrorResponse(errcode.ErrorImportCsvFile)
		return
	}
	processInfo, err := a_procession.NewProcessionService(c, global.DBEngine, global.CacheEngine).ProcessAll()
	if err != nil {
		global.Logger.Errorf(c, "svc.ProcessAll err: %v", err)
		response.ToErrorResponse(errcode.ErrorProcess)
		return
	}

	response.ToResponse(gin.H{
		"ImportInfos": importDataInfos,
		"Process":     processInfo,
	})
}
