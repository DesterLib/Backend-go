package routes

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name   string
	Engine *gin.Engine
}

func (r *Route) Init(engine *gin.Engine) {
	r.Engine = engine
}

func NewRoute(name string) *Route {
	return &Route{Name: strings.TrimSuffix(name, "/")}
}

func (r *Route) GET(s string, handlers ...gin.HandlerFunc) {
	r.Engine.GET(fmt.Sprintf("%s/%s", r.Name, s), handlers...)
}

func (r *Route) POST(s string, handlers ...gin.HandlerFunc) {
	r.Engine.POST(fmt.Sprintf("%s/%s", r.Name, s), handlers...)
}
