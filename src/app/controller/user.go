package controller

import (
	"errors"
	"encoding/json"
	"net/http"
	"app/model"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

type Success struct {
	Success string `json: "success"`
	User []model.User `json: "user"`
}

type Error struct {
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
			if err.Tag() == "alphanum"|| err.Tag() == "excludesall" {
				error += err.Field() + ": constains invalid characters "
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
		response := &Error {
			Err: e,
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	} else if ex != nil {
		s := ex.Error()
		response := &Error{
			Err: s,
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	} else {
		response := &Success{
			Success: "user created",
			User: []model.User{user}}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
