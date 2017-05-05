package controller

import (
	"app/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// struct para gerar a hash do JWT
type Claims struct {
	User model.User `json:"user"`
	jwt.StandardClaims
}

// struct de resposta de sucessso
type successLogin struct {
	Success string `json:"success"`
	Token   string `json:"token"`
}

// struct de resposta de fracasso
type errorLogin struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

const (
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

// handler de autenticacao do usuario
// compara a senha inserida com a senha do banco
// verifica se o usuario existe no banco
// retorna o hash de autenticacao ou um mensagem de erro

func AuthPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	name, password := r.FormValue("name"), r.FormValue("password")
	user, gErr := model.UserByName(name)
	var jData []byte
	if gErr != nil {
		response := &errorLogin{
			Error: gErr.Error(),
			Token: ""}
		jData, _ = json.Marshal(response)
	} else if user.Password == password {
		claims := Claims{
			user,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
				Issuer:    "localhost:3000"}}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
		tokenString, tErr := token.SignedString([]byte(mySigningKey))
		if tErr != nil {
			response := &errorLogin{
				Error: tErr.Error(),
				Token: ""}
			jData, _ = json.Marshal(response)
		} else {
			response := &successLogin{
				Success: "user authenticated",
				Token:   tokenString}
			jData, _ = json.Marshal(response)
		}
	} else {
		response := &errorLogin{
			Error: "invalid password",
			Token: ""}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}

// middleware de validacao simples para o JWT
// verifica a validade do token armazenado no cabecalho
// retorna erro ou deixa o usuario acessar a url desejada

func Validate(protectedPage httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var jData []byte
		xtoken := r.Header.Get("X-Token")
		if len(xtoken) <= 0 {
			response := &errorLogin{
				Error: "user not validated",
				Token: xtoken}
			jData, _ = json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		token, err := jwt.ParseWithClaims(xtoken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected siging method")
			}
			return []byte(mySigningKey), nil
		})
		if err != nil {
			response := &errorLogin{
				Error: err.Error(),
				Token: xtoken}
			jData, _ = json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), mySigningKey, *claims)
			protectedPage(w, r.WithContext(ctx), p)
		} else {
			http.NotFound(w, r)
			return
		}
	}
}
