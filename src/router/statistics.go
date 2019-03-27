package router

import (
	"controller"
	"github.com/gin-gonic/gin"
)

func RouteStatistics(g *gin.RouterGroup) {
	g.GET("/info", controller.GetStatsInfo)

	g.GET("/id", controller.GetRecordsByID)
	g.GET("/formula", controller.GetRecordsByFormula)
}
