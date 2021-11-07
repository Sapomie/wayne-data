package routers

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/Sapomie/wayne-data/internal/service/essential"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/test", func(c *gin.Context) {
		//_, info, err := rawEvent.ImportCsvData()
		//if err != nil {
		//	fmt.Println(err)
		//}

		start, _ := mtime.NewTimeZone(mtime.TypeYear, 2021, 1).BeginAndEnd()
		events, _, err := model.NewEventModel(global.DBEngine).GetAll()
		ess, err := essential.MakeEssentials(events, start, cons.Newest, mtime.TypeMonth)
		if err != nil {
			return
		}

		c.JSON(200, ess)
	})

	apiv1 := r.Group("/api/v1")

	{
		apiv1.GET("event", v1.ListEvents)

	}

	return r
}
