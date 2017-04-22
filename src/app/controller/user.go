package controller

import (
	"net/http"
	"app/model"
	"log"
	"github.com/julienschmidt/httprouter"
)

func UserPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, password := r.FormValue("name"), r.FormValue("password")
	ex := model.UserCreate(name, password)
	if ex != nil {
		log.Println(ex)
	}
}
