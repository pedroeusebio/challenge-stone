package route

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func LoadHTTP() http.Handler {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()

	return r
}
