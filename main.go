package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/api"
	"github.com/desterlib/backend-go/cache"
	"github.com/desterlib/backend-go/config"
	"github.com/desterlib/backend-go/db"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New(os.Stderr, &logger.LoggerOpts{
		ProjectName:  "DESTER",
		MinimumLevel: logger.LevelInfo,
	})
	log.Println("STARTING...")
	log.ChangeLevel(logger.LevelInfo)
	config.Load(log)
	cache.Load(log)
	db.LoadDB(log)
	router := gin1(log)
	log.Printlnf("STARTED AT PORT: '%d'\n", config.ValueOf.Port)
	router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
}

func gin1(l *logger.Logger) *gin.Engine {
	// log := l.Create("GIN")
	if config.ValueOf.DesterDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "index.html", gin.H{})
	})
	router.Static("static", "build/static")
	router.StaticFile("favicon.ico", "./build/favicon.ico")
	router.StaticFile("asset-manifest.json", "./build/asset-manifest.json")
	router.StaticFile("/", "./build/index.html")
	router.LoadHTMLFiles("build/index.html")
	api.Load(router, l)
	return router
}
