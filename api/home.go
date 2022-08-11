package api

import (
	"net/http"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/db"
	"github.com/desterlib/backend-go/routes"
	"github.com/desterlib/backend-go/types"
	"github.com/gin-gonic/gin"
)

func (e *entry) LoadHome(r *routes.Route) {
	log := e.Logger.Create("HOME")
	log.ChangeLevel(logger.LevelInfo)
	defer log.Println("LOADED ROUTE")
	r.GET("home", getHome)
}

func getHome(ctx *gin.Context) {
	config := db.GetConfig()
	if config.App.SecretKey == "" {
		ctx.JSON(http.StatusPreconditionRequired, types.DataResponse{
			Code:        http.StatusPreconditionRequired,
			Message:     "The config needs to be initialized first.",
			Ok:          false,
			Result:      "/settings",
			TimeTaken:   0,
			Title:       config.App.Title,
			Description: config.App.Description,
		})
	} else {
		ctx.JSON(http.StatusOK, types.DataResponse{
			Code:        http.StatusOK,
			Message:     "Home page data successfully retrieved.",
			Ok:          true,
			Result:      "",
			Title:       config.App.Title,
			Description: config.App.Description,
		})
	}
}
