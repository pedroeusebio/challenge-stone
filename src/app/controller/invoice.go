package controller

import (
	"strconv"
	"net/http"
	"app/model"
	"app/shared/ordenate"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/go-playground/validator.v9"
)

type successInvoice struct {
	Success []string `json:"success"`
	Invoice []model.Invoice `json:"invoice"`
}

type errorInvoice struct {
	Err []string `json:"error"`
	Invoice []model.Invoice `json:"invoice"`
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
	invoice := model.Invoice{Amount: amount,Document: document,Month: month,Year: year,Is_active: true}
	vErr := validateInvoice(invoice)
	var jData []byte
	if err1 != nil || err2 != nil || err3 != nil {
		response := &errorInvoice{
			Err: []string{"error while parsing values"},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else if len(vErr) > 0 {
		response := &errorInvoice {
			Err: vErr,
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else {
		ex := model.InvoiceCreate(amount, document, month, year)
		if ex != nil {
			s := ex.Error()
			response := &errorInvoice{
				Err: []string{s},
				Invoice: []model.Invoice{invoice}}
			jData, _ = json.Marshal(response)
		} else {
			response := &successInvoice{
				Success: []string{"invoice_create"},
				Invoice: []model.Invoice{invoice}}
			jData, _ = json.Marshal(response)
		}
}
	w.Write(jData)
}

func InvoiceGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	orderBy, page, length := r.FormValue("order"), r.FormValue("page"), r.FormValue("length")
	var order []ordenate.Ordenate
	var oErr error
	if len(orderBy) > 0 {
		order, oErr = ordenate.Order(orderBy)
	} else {
		order, oErr = []ordenate.Ordenate{}, nil
	}
	users, err := model.InvoiceGetAll(order, page, length)
	var jData []byte
	if oErr != nil {
		response := &errorInvoice {
			Err: []string{oErr.Error()},
			Invoice: []model.Invoice{}}
		jData, _ = json.Marshal(response)
	} else if err != nil {
		response := &errorInvoice {
			Err: []string{err.Error()},
			Invoice: []model.Invoice{}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successInvoice {
			Success: []string{"invoice_getall"},
			Invoice: users}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}

func InvoiceDEL(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	invoice, ex := model.InvoiceDelete(id)
	invoice.Is_active = false
	var jData []byte
	if ex != nil {
		s := ex.Error()
		response := &errorInvoice {
			Err: []string{s},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successInvoice {
			Success: []string{"invoice_delete"},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
