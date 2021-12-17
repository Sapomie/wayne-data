package app

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"github.com/Sapomie/wayne-data/pkg/convert"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("page_size")).MustInt()
	if pageSize <= 0 {
		return global.AppSetting.DefaultPageSize
	}
	if pageSize > global.AppSetting.MaxPageSize {
		return global.AppSetting.MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}

func GetLimitOffset(c *gin.Context) (limit, offset int) {
	limit = GetPageSize(c)
	offset = GetPageOffset(GetPage(c), limit)
	return
}

func GinBeginAndEnd(c *gin.Context) (start, end time.Time) {

	var (
		date    = c.DefaultQuery("date", "")
		spanStr = c.DefaultQuery("span", "")
	)
	span, _ := strconv.Atoi(spanStr)

	start, end, err := DateStartAndEnd(date, span)
	if err != nil {
		return
	}

	return
}

func DateStartAndEnd(date string, span int) (start, end time.Time, err error) {
	start = cons.Oldest
	end = cons.Futurest
	if len(date) < 8 {
		//err = errors.New("wrong date format")
		return
	}
	yearStr := date[:4]
	monthStr := date[4:6]
	dayStr := date[6:8]

	year, err := strconv.Atoi(yearStr)
	month, err := strconv.Atoi(monthStr)
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return
	}

	start = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	end = time.Date(year, time.Month(month), day+1+span, 0, 0, 0, 0, time.Local)
	return
}
