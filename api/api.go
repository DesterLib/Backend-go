package api

import (
	"reflect"

	"github.com/anonyindian/logger"
	"github.com/desterlib/backend-go/config"
	"github.com/desterlib/backend-go/routes"
	"github.com/gin-gonic/gin"
)

type entry struct {
	Logger *logger.Logger
}

func Load(r *gin.Engine, l *logger.Logger) {
	l = l.Create("API")
	l.ChangeLevel(logger.LevelMain)
	defer l.Println("LOADED ALL API ROUTES")
	route := routes.NewRoute(config.DEFAULT_API_V1_ENDPOINT)
	route.Init(r)
	Type := reflect.TypeOf(&entry{l})
	Value := reflect.ValueOf(&entry{l})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(route)})
	}
}
