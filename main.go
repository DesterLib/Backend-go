package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/api"
	"github.com/desterlib/backend-go/config"
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
	ginl()
	log.Printlnf("STARTED AT PORT: '%d'", config.ValueOf.Port)
}

func ginl() {
	router := gin.Default()

	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "index.html", gin.H{})
	})
	router.Static("static", "build/static")
	router.StaticFile("favicon.ico", "./build/favicon.ico")
	router.StaticFile("asset-manifest.json", "./build/asset-manifest.json")
	router.LoadHTMLFiles("build/index.html")
	api.Load(router)
	router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))

}
