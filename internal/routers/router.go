package routers

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/routers/v1"
	"github.com/Sapomie/wayne-data/internal/service/c_book"
	"github.com/Sapomie/wayne-data/internal/service/c_series"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("view/*.html")

	r.GET("/test", func(c *gin.Context) {
		//_, info, err := b_rawEvent.NewRawEventService(c, global.DBEngine, global.CacheEngine).ImportCsvData()
		//if err != nil {
		//	fmt.Println(err)
		//}
		infos, err := c_book.NewBookService(c, global.DBEngine, global.CacheEngine).ProcessBook()
		if err != nil {
			fmt.Println(err)
		}
		infos, err = c_series.NewSeriesService(c, global.DBEngine, global.CacheEngine).ProcessSeries()
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, infos)
	})

	apiv1 := r.Group("/api/v1")

	{
		apiv1.Static("static/", "view/static")
		apiv1.GET("event", v1.ListEvents)

		apiv1.GET("essday", v1.ListEssentialsDay)
		apiv1.GET("essten", v1.ListEssentialsTen)
		apiv1.GET("essmonth", v1.ListEssentialsMonth)
		apiv1.GET("essquarter", v1.ListEssentialsQuarter)
		apiv1.GET("essyear", v1.ListEssentialsYear)

		apiv1.GET("progressnow/:typ", v1.GetProgressMonthNow)
	}

	return r
}
