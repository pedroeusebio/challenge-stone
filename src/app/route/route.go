package route

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"app/controller"
)

func LoadHTTP() http.Handler {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()

	r.POST("/user", controller.UserPOST)
	r.GET("/user", controller.UserGET)
	return r
}
