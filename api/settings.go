package api

import (
	"io/ioutil"
	"net/http"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/db"
	"github.com/desterlib/backend-go/rclone"
	"github.com/desterlib/backend-go/routes"
	"github.com/desterlib/backend-go/types"
	"github.com/gin-gonic/gin"
)

func (e *entry) LoadSettings(r *routes.Route) {
	log := e.Logger.Create("SETTINGS")
	log.ChangeLevel(logger.LevelInfo)
	defer log.Println("LOADED ROUTE")
	r.GET("settings", settingsGet)
	r.POST("settings", settingsPost)
}

func settingsPost(ctx *gin.Context) {
	secret_key := ctx.Query("secret_key")
	config := db.GetConfig()
	if secret_key == config.App.SecretKey {
		b, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}
		go db.SaveConfig(b)
		rclone.Restart()
	} else {
		ctx.JSON(http.StatusOK, types.DataResponse{
			Code:        http.StatusUnauthorized,
			Message:     "The secret key was incorrect.",
			Ok:          false,
			Result:      nil,
			TimeTaken:   0,
			Title:       "Dester",
			Description: "Dester",
		})
	}
}

func settingsGet(ctx *gin.Context) {
	secret_key := ctx.Query("secret_key")
	config := db.GetConfig()
	if secret_key == config.App.SecretKey {
		ctx.JSON(http.StatusOK, types.DataResponse{
			Code:        http.StatusOK,
			Message:     "Config successfully retrieved from database.",
			Ok:          true,
			Result:      config,
			TimeTaken:   0,
			Title:       "Dester",
			Description: "Dester",
		})
	} else {
		ctx.JSON(http.StatusOK, types.DataResponse{
			Code:        http.StatusUnauthorized,
			Message:     "The secret key was incorrect.",
			Ok:          false,
			Result:      nil,
			TimeTaken:   0,
			Title:       "Dester",
			Description: "Dester",
		})
	}
}
