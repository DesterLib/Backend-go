package api

import (
	"net/http"

	"github.com/desterlib/backend-go/routes"
	"github.com/desterlib/backend-go/utils"
	"github.com/gin-gonic/gin"
)

type wakanda struct {
	Auth0      struct{} `json:"auth0"`
	Categories []string `json:"categories"`
	Gdrive     struct{} `json:"gdrive"`
	Onedrive   struct{} `json:"onedrive"`
	Sharepoint struct{} `json:"sharepoint"`
	Tmdb       struct{} `json:"tmdb"`
	Subtitles  struct{} `json:"subtitles"`
	Build      struct{} `json:"build"`
	Rclone     []string `json:"rclone"`
}

func (*entry) LoadSettings(r *routes.Route) {
	r.GET("settings", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.Response{
			Code:        http.StatusOK,
			Message:     "Config successfully retrieved from database.",
			Ok:          true,
			Result:      wakanda{},
			TimeTaken:   0,
			Title:       "Dester",
			Description: "Dester",
		})
	})
}
