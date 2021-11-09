package v1

import (
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/service/essential"
	"github.com/Sapomie/wayne-data/pkg/app"
	"github.com/Sapomie/wayne-data/pkg/errcode"
	"github.com/Sapomie/wayne-data/pkg/mtime"
	"github.com/gin-gonic/gin"
)

func ListEssentialsDay(c *gin.Context) {
	response := app.NewResponse(c)
	svc := essential.NewEssentialService(c)
	ess, _, err := svc.GetEssentialList(mtime.TypeDay)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}
	response.ToResponseHtml("essential.html", gin.H{
		"resp":      ess,
		"type":      "day",
		"tableName": "datatableDay",
	})
}

func ListEssentialsTen(c *gin.Context) {
	response := app.NewResponse(c)
	svc := essential.NewEssentialService(c)
	ess, _, err := svc.GetEssentialList(mtime.TypeTen)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}
	response.ToResponseHtml("essential.html", gin.H{
		"resp":      ess,
		"type":      "ten",
		"tableName": "datatableTen",
	})
}

func ListEssentialsMonth(c *gin.Context) {
	response := app.NewResponse(c)
	svc := essential.NewEssentialService(c)
	ess, _, err := svc.GetEssentialList(mtime.TypeMonth)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}

	response.ToResponseHtml("essential.html", gin.H{
		"resp":      ess,
		"type":      "month",
		"tableName": "datatableMonth",
	})
}

func ListEssentialsQuarter(c *gin.Context) {
	response := app.NewResponse(c)
	svc := essential.NewEssentialService(c)
	ess, _, err := svc.GetEssentialList(mtime.TypeQuarter)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}

	response.ToResponseHtml("essential.html", gin.H{
		"resp":      ess,
		"type":      "quarter",
		"tableName": "datatableQuarter",
	})
}

func ListEssentialsYear(c *gin.Context) {
	response := app.NewResponse(c)
	svc := essential.NewEssentialService(c)
	essHalf, _, err := svc.GetEssentialList(mtime.TypeHalf)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}
	essYear, _, err := svc.GetEssentialList(mtime.TypeYear)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetEssentialList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetEssentialListFail)
		return
	}

	essHalf = append(essHalf, essYear...)

	response.ToResponseHtml("essential.html", gin.H{
		"resp":      essHalf,
		"type":      "year",
		"tableName": "datatableYear",
	})
}
