package routers

import (
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "ok")
	})

	apiv1 := r.Group("/api/v1")

	{
		apiv1.GET("event", v1.ListEvents)

	}

	return r
}
