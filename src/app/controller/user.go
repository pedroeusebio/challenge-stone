package controller

import (
	"errors"
	"encoding/json"
	"net/http"
	"app/model"
	"app/shared/ordenate"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

type successUser struct {
	Success string `json: "success"`
	User []model.User `json: "user"`
}

type errorUser struct {
	Err string `json: "error"`
	User []model.User `json: "user"`
}

var validate *validator.Validate


func validateUser(user model.User) error {
	var error string
	vErr := validate.Struct(user)
	if vErr != nil {
		for _, err := range vErr.(validator.ValidationErrors) {
			error += ", "
			if err.Tag() == "required" {
				error += err.Field() + ": is required "
			}
			if err.Tag() == "alphanum" || err.Tag() == "excludesall" {
				error += err.Field() + ": contains invalid characters "
			}
			if err.Tag() == "gt" {
				var gt string
				if err.Field() == "Name" {
					gt = model.GtName
				} else {
					gt = model.GtPassword
				}
				error += err.Field() + ": must have more than " + gt + " characters"
			}
		}
		rErr := errors.New(error)
		return rErr
	} else {
		return nil
	}
}

func UserPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	validate = validator.New()
	w.Header().Set("Content-Type", "application/json")
	name, password := r.FormValue("name"), r.FormValue("password")
	user := model.User{name, password}
	vErr := validateUser(user)
	ex := model.UserCreate(name, password)
	var jData []byte
	if vErr != nil {
		e := vErr.Error()
		response := &errorUser {
			Err: e,
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	} else if ex != nil {
		s := ex.Error()
		response := &errorUser {
			Err: s,
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successUser {
			Success: "user_create",
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}


func UserGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	orderBy, page, length := r.FormValue("order"), r.FormValue("page"), r.FormValue("length")
	order, oErr := ordenate.Order(orderBy)
	users, err := model.UserGetAll(order, page, length)
	var jData []byte
	if oErr != nil {
		response := &errorUser {
			Err: oErr.Error(),
			User: []model.User{}}
		jData, _ = json.Marshal(response)
	} else if err != nil {
		response := &errorUser {
			Err: err.Error(),
			User: []model.User{}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successUser {
			Success: "user_getall",
			User: users}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
