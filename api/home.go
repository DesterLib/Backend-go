package api

import (
	"net/http"

	"github.com/desterlib/backend-go/routes"
	"github.com/desterlib/backend-go/utils"
	"github.com/gin-gonic/gin"
)

func (*entry) LoadHome(r *routes.Route) {
	r.GET("home", func(ctx *gin.Context) {
		ctx.JSON(http.StatusPreconditionRequired, utils.Response{
			Code:        http.StatusPreconditionRequired,
			Message:     "The config needs to be initialized first.",
			Ok:          false,
			Result:      "/settings",
			TimeTaken:   0,
			Title:       "Dester",
			Description: "Dester",
		})
	})
}
