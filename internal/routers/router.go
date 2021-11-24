package routers

import (
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("view/*.html")

	router.GET("/test", func(c *gin.Context) {

		c.JSON(200, "haha")
	})

	apiv1 := router.Group("/api/v1")

	{
		apiv1.Static("static/", "view/static")

		apiv1.GET("essday", v1.ListEssentialsDay)
		apiv1.GET("essten", v1.ListEssentialsTen)
		apiv1.GET("essmonth", v1.ListEssentialsMonth)
		apiv1.GET("essquarter", v1.ListEssentialsQuarter)
		apiv1.GET("essyear", v1.ListEssentialsYear)

		apiv1.GET("progressnow/:typ", v1.GetProgressNow)

		apiv1.GET("upload", v1.Upload)
		apiv1.POST("upload", v1.UploadPost)

		apiv1.GET("event", v1.ListEvents)

		apiv1.GET("export", v1.ExportAllRawEvent)
	}

	return router
}
