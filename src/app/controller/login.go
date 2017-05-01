package controller

import (
	"time"
	"encoding/json"
	"net/http"
	"app/model"
	"github.com/julienschmidt/httprouter"
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type successLogin struct {
	Success string `json:"success"`
	Token string `json:"token"`
}

type errorLogin struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

const (
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

func AuthPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	name, password := r.FormValue("name"), r.FormValue("password")
	user, gErr := model.UserByName(name)
	var jData []byte
	if gErr != nil {
		response := &errorLogin {
			Error: gErr.Error(),
			Token: ""}
		jData, _ = json.Marshal(response)
	} else if user.Password == password {
		claims := Claims {
			user.Name + user.Password,
			jwt.StandardClaims {
				ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
				Issuer: "localhost:3000"}}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"),claims)
		tokenString, tErr := token.SignedString([]byte(mySigningKey))
		if tErr != nil {
			response := &errorLogin {
				Error: tErr.Error(),
				Token: ""}
			jData, _ = json.Marshal(response)
		} else {
			response := &successLogin {
				Success: "user authenticated",
				Token: tokenString}
			jData, _ = json.Marshal(response)
		}
	} else {
		response := &errorLogin {
			Error: "invalid password",
			Token: ""}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
