package route

import (
	"app/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func LoadHTTP() http.Handler {
	return routes()
}

func routes() *httprouter.Router {
	r := httprouter.New()
	// rotas do user
	r.GET("/user", controller.UserGET)
	r.POST("/user", controller.UserPOST)
	// rotas do invoice
	// controller.Validate aplica a restricao da rota para somente os autenticados
	r.DELETE("/invoice/:id", controller.Validate(controller.InvoiceDEL))
	r.GET("/invoice", controller.InvoiceGET)
	r.POST("/invoice", controller.Validate(controller.InvoicePOST))
	// rota de login
	r.POST("/login", controller.AuthPOST)
	return r
}
