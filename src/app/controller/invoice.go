package controller

import (
	"app/model"
	"app/shared/ordenate"
	v "app/shared/validate"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// struct de resposta de sucesso
type successInvoice struct {
	Success []string        `json:"success"`
	Invoice []model.Invoice `json:"payload"`
}

// struct de resposta de erro
type errorInvoice struct {
	Err     []string        `json:"error"`
	Invoice []model.Invoice `json:"payload"`
}

// handler de criacao do invoice.
// os dados são validados e salvos no banco.
// retorna um json com mensagem de sucesso ou erro e um payload dos dados do invoice

func InvoicePOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	amount, err1 := strconv.ParseFloat(r.FormValue("amount"), 64)
	document := r.FormValue("document")
	month, err2 := strconv.Atoi(r.FormValue("month"))
	year, err3 := strconv.Atoi(r.FormValue("year"))
	invoice := model.Invoice{Amount: amount, Document: document, Month: month, Year: year, Is_active: true}
	vErr := v.ValidateInvoice(invoice)
	var jData []byte
	if err1 != nil || err2 != nil || err3 != nil {
		response := &errorInvoice{
			Err:     []string{"error while parsing values"},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else if len(vErr) > 0 {
		response := &errorInvoice{
			Err:     vErr,
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else {
		ex := model.InvoiceCreate(amount, document, month, year)
		if ex != nil {
			s := ex.Error()
			response := &errorInvoice{
				Err:     []string{s},
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

// handler da listagem do invoice.
// realiza a busca no banco e retorna um array de invoices ou um array vazio e uma mensagem de erro.

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
		response := &errorInvoice{
			Err:     []string{oErr.Error()},
			Invoice: []model.Invoice{}}
		jData, _ = json.Marshal(response)
	} else if err != nil {
		response := &errorInvoice{
			Err:     []string{err.Error()},
			Invoice: []model.Invoice{}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successInvoice{
			Success: []string{"invoice_getall"},
			Invoice: users}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}

// handler da remocao do invoice_.
// recebe o id do invoice e o remove do banco
// retorna sucesso e o invoice removido ou uma mensagem de erro e um invoice vazio.
func InvoiceDEL(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := params.ByName("id")
	invoice, ex := model.InvoiceDelete(id)
	invoice.Is_active = false
	var jData []byte
	if ex != nil {
		s := ex.Error()
		response := &errorInvoice{
			Err:     []string{s},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	} else {
		response := &successInvoice{
			Success: []string{"invoice_delete"},
			Invoice: []model.Invoice{invoice}}
		jData, _ = json.Marshal(response)
	}
	w.Write(jData)
}
