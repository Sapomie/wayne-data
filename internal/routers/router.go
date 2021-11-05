package routers

import (
	"fmt"
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/Sapomie/wayne-data/internal/service/rawEvent"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/test", func(c *gin.Context) {
		_, info, err := rawEvent.ImportCsvData()
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, info)
	})

	apiv1 := r.Group("/api/v1")

	{
		apiv1.GET("event", v1.ListEvents)

	}

	return r
}
