package routers

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/Sapomie/wayne-data/internal/service/b_rawEvent"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("view/*.html")

	router.GET("/test", func(c *gin.Context) {

		err := b_rawEvent.NewRawEventService(c, global.DBEngine, global.CacheEngine).ReadDefaultTaskValue()
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, "haha")
	})

	apiv1 := router.Group("/api/v1")

	{
		apiv1.Static("static/", "view/static")

		//essential
		apiv1.GET("essday", v1.ListEssentialsDay)
		apiv1.GET("essten", v1.ListEssentialsTen)
		apiv1.GET("essmonth", v1.ListEssentialsMonth)
		apiv1.GET("essquarter", v1.ListEssentialsQuarter)
		apiv1.GET("essyear", v1.ListEssentialsYear)

		//progress
		apiv1.GET("progressnow/:typ", v1.GetProgressNow)

		//event
		apiv1.GET("event", v1.ListEvents)

		//upload
		apiv1.GET("upload", v1.Upload)
		apiv1.POST("upload", v1.UploadPost)
		apiv1.GET("export", v1.ExportAllRawEvent)

		//mobile
		apiv1.GET("mbprogressnow/:typ", v1.GetMobileProgressNow)
	}

	return router
}
