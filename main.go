package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/desterlib/backend-go/types"
	"github.com/desterlib/backend-go/utils"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/api"
	"github.com/desterlib/backend-go/cache"
	"github.com/desterlib/backend-go/config"
	"github.com/desterlib/backend-go/db"
	"github.com/gin-gonic/gin"
)

var versionString string = "0.0.1"

var startTime time.Time = time.Now()

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
	log.Printlnf("SERVER STARTED AT: http://localhost:%d\n", config.ValueOf.Port)
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
	// set up the middleware for cors
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	})
	if _, err := os.Stat("build"); err == nil {
		router.NoRoute(func(ctx *gin.Context) {
			ctx.HTML(http.StatusNotFound, "index.html", gin.H{})
		})
		router.Static("static", "build/static")
		router.StaticFile("favicon.ico", "./build/favicon.ico")
		router.StaticFile("asset-manifest.json", "./build/asset-manifest.json")
		router.StaticFile("/", "./build/index.html")
		router.LoadHTMLFiles("build/index.html")
	} else if errors.Is(err, os.ErrNotExist) {
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, types.RootResponse{
				Message: "Backend is working.",
				Ok:      true,
				Uptime:  utils.TimeFormat(uint64(time.Since(startTime).Seconds())),
				Version: versionString,
			})
		})
	} else {
		panic(err)
	}
	api.Load(router, l)
	return router
}
