package main

import (
	"config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"router"
)

func main() {
	startWeb()
}

func startWeb() {
	// TODO: Log
	// TODO: 在这里建立redis索引
	// TODO: 是不是可以直接缓存每个页面

	e := gin.Default()

	// Cors - allow all origins
	e.Use(cors.Default())

	// Router
	g := e.Group("/stats")
	router.RouteStatistics(g)

	log.Fatal(e.Run(config.GlobalConfig.Addr))
}
