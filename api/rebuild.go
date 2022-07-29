package api

import (
	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/routes"
	"github.com/gin-gonic/gin"
)

func (e *entry) LoadRebuild(r *routes.Route) {
	log := e.Logger.Create("REBUILD")
	log.ChangeLevel(logger.LevelInfo)
	defer log.Println("LOADED ROUTE")
	r.GET("rebuild", getRebuild)
}

func getRebuild(ctx *gin.Context) {
}
