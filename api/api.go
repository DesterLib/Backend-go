package api

import (
	"reflect"

	"github.com/desterlib/backend-go/routes"
	"github.com/gin-gonic/gin"
)

type entry struct{}

func Load(r *gin.Engine) {
	route := routes.NewRoute("/api/v1")
	route.Init(r)
	Type := reflect.TypeOf(&entry{})
	Value := reflect.ValueOf(&entry{})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(route)})
	}
}
