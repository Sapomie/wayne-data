package routers

import (
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("view/*.html")

	r.GET("/test", func(c *gin.Context) {

		c.JSON(200, "haha")
	})

	apiv1 := r.Group("/api/v1")

	{
		apiv1.Static("static/", "view/static")

		apiv1.GET("essday", v1.ListEssentialsDay)
		apiv1.GET("essten", v1.ListEssentialsTen)
		apiv1.GET("essmonth", v1.ListEssentialsMonth)
		apiv1.GET("essquarter", v1.ListEssentialsQuarter)
		apiv1.GET("essyear", v1.ListEssentialsYear)

		apiv1.GET("progressnow/:typ", v1.GetProgressMonthNow)

		apiv1.GET("upload", v1.Upload)
		apiv1.POST("upload", v1.UploadPost)

		apiv1.GET("event", v1.ListEvents)
	}

	return r
}
