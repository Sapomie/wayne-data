package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/c_movie"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListMovies(c *gin.Context) {
	response := app.NewResponse(c)

	svc := c_movie.NewMovieService(c, global.DBEngine, global.CacheEngine)
	resp, err := svc.ListMovies()
	if err != nil {
		global.Logger.Errorf(c, "svc.ListMovies err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetMovies)
		return
	}

	response.ToResponseHtml("movie.html", gin.H{
		"resp": resp.Items,
		"sum":  resp.Sum,
	})
}
