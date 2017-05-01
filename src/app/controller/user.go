package controller

import (
	"app/model"
	"app/shared/ordenate"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

type successUser struct {
	Success string       `json:"success"`
	User    []model.User `json:"payload"`
}

type errorUser struct {
	Err  string       `json:"error"`
	User []model.User `json:"payload"`
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
	var jData []byte
	if vErr != nil {
		e := vErr.Error()
		response := &errorUser{
			Err:  e,
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	} else {
		ex := model.UserCreate(name, password)
		if ex != nil {
			s := ex.Error()
			response := &errorUser{
				Err:  s,
				User: []model.User{user}}
			jData, _ = json.Marshal(response)
		} else {
			response := &successUser{
				Success: "user_create",
				User:    []model.User{user}}
			jData, _ = json.Marshal(response)
		}
	}
	w.Write(jData)
}

func UserGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	orderBy, page, length, name := r.FormValue("order"), r.FormValue("page"), r.FormValue("length"), r.FormValue("name")
	var order []ordenate.Ordenate
	var oErr error
	if len(orderBy) > 0 {
		order, oErr = ordenate.Order(orderBy)
	} else {
		order, oErr = []ordenate.Ordenate{}, nil
	}
	users, err := model.UserGetAll(order, page, length, name)
	var jData []byte
	if oErr != nil {
		response := &errorUser{
			Err:  oErr.Error(),
			User: []model.User{}}
		jData, _ = json.Marshal(response)
	} else if err != nil {
		response := &errorUser{
			Err:  err.Error(),
			User: []model.User{}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successUser{
			Success: "user_getall",
			User:    users}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
