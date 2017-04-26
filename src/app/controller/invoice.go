package controller

import (
	"strconv"
	"net/http"
	"app/model"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

type successInvoice struct {
	Success []string `json: "success"`
	User []model.Invoice `json: "invoice"`
}

type errorInvoice struct {
	Err []string `json: "error"`
	User []model.Invoice `json: "invoice"`
}



func validateInvoice(invoice model.Invoice) []string {
	error := []string{}
	Err := validate.Struct(invoice)
	if Err != nil {
		for _, err := range Err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				error = append(error, err.Field() + ": is required ")
			}
			if err.Tag() == "gte" {
				var gte string
				if err.Field() == "Amount" {gte = model.GteAmount} else {gte = model.GteMonth}
				error = append(error, err.Field() + ": must be greater than or equals to " + gte + " ")
			}
			if err.Tag() == "gt" {
				error = append(error, err.Field() + ": must be greater than " + model.GtYear + " ")
			}
			if err.Tag() == "lte" {
				error = append(error, err.Field() + ": must be less than or equals to " + model.LteMonth + " ")
			}
		}
	}
	return error
}

func InvoicePOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	validate = validator.New()
	w.Header().Set("Content-Type", "application/json")
	amount, err1 := strconv.ParseFloat(r.FormValue("amount"), 64)
	document := r.FormValue("document")
	month, err2 := strconv.Atoi(r.FormValue("month"))
	year, err3 := strconv.Atoi(r.FormValue("year"))
	invoice := model.Invoice{amount, document, month, year, true}
	vErr := validateInvoice(invoice)
	ex := model.InvoiceCreate(amount, document, month, year)
	var jData []byte
	if err1 != nil || err2 != nil || err3 != nil {
		response := &errorInvoice{
			Err: []string{"error while parsing values"},
			User: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else if len(vErr) > 0 {
		response := &errorInvoice {
			Err: vErr,
			User: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else if ex != nil {
		s := ex.Error()
		response := &errorInvoice{
			Err: []string{s},
			User: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successInvoice{
			Success: []string{"invoice_create"},
			User: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
