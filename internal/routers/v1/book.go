package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/c_book"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func ListBooks(c *gin.Context) {
	response := app.NewResponse(c)

	svc := c_book.NewBookService(c, global.DBEngine, global.CacheEngine)
	books, sum, err := svc.ListBooks()
	if err != nil {
		global.Logger.Errorf(c, "svc.NewBookService err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetBook)
		return
	}

	response.ToResponseHtml("book.html", gin.H{
		"resp": books,
		"sum":  sum,
	})
}
