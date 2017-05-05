package controller

import (
	"app/shared/database"
	"app/shared/server"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type configuration struct {
	Database database.Database `json:"Database"`
	Server   server.Server     `json:"Server"`
}

func ParseJsonFile(configPath string) configuration {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("File error %v", err)
		os.Exit(1)
	}
	var config configuration
	json.Unmarshal(file, &config)
	return config
}
func init() {
	config := ParseJsonFile("../../../config/config.json")
	database.Connect(config.Database)
}

func TestInvoicePOST1(t *testing.T) {
	invoice := url.Values{}
	invoice.Set("amount", "123")
	invoice.Add("document", "11004735480")
	invoice.Add("month", "2")
	invoice.Add("year", "12")
	request, _ := http.NewRequest("POST", "/invoice", bytes.NewBufferString(invoice.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	InvoicePOST(response, request, nil)
	result := successInvoice{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, "invoice_create", result.Success[0], "Ok response is expected")
}

func TestInvoicePOST2(t *testing.T) {
	invoice := url.Values{}
	invoice.Set("amount", "123")
	invoice.Add("document", "11004735481") //invalid cpf
	invoice.Add("month", "2")
	invoice.Add("year", "12")
	request, _ := http.NewRequest("POST", "/invoice", bytes.NewBufferString(invoice.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	InvoicePOST(response, request, nil)
	result := errorInvoice{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, "Document: CPF invalid", result.Err[0], "should return error")
}

//test geting the first page with 5 invoices
func TestInvoiceGET1(t *testing.T) {
	invoice := url.Values{}
	invoice.Set("page", "0")
	invoice.Add("length", "5")
	request, _ := http.NewRequest("GET", "/invoice?"+invoice.Encode(), nil)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	InvoiceGET(response, request, nil)
	result := successInvoice{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, []string{"invoice_getall"}, result.Success, "should return error")
	assert.True(t, 5 >= len(result.Invoice), "should have the less or equal than 5")
}

//test wrong column
func TestInvoiceGET2(t *testing.T) {
	invoice := url.Values{}
	invoice.Set("page", "0")
	invoice.Add("order", "[{\"Column\":\"asdasd\",\"Order\":\"desc\"}]")
	request, _ := http.NewRequest("GET", "/invoice?"+invoice.Encode(), nil)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	InvoiceGET(response, request, nil)
	result := errorInvoice{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, []string{"pq: column \"asdasd\" does not exist"}, result.Err, "should return error")
}

func TestInvoiceDEL1(t *testing.T) {
	request, _ := http.NewRequest("DEL", "/invoice/", nil)
	response := httptest.NewRecorder()
	params := httprouter.Param{
		Key:   "id",
		Value: "2"}
	InvoiceDEL(response, request, httprouter.Params{params})
	result := successInvoice{}
	error := errorInvoice{}
	json.NewDecoder(response.Body).Decode(&result)
	json.NewDecoder(response.Body).Decode(&error)
	assert.Equal(t, []string{"invoice_delete"}, result.Success, "should return invoice_delete") //if exists
}

func TestInvoiceDEL2(t *testing.T) {
	request, _ := http.NewRequest("DEL", "/invoice/", nil)
	response := httptest.NewRecorder()
	params := httprouter.Param{
		Key:   "id",
		Value: "999999"}
	InvoiceDEL(response, request, httprouter.Params{params})
	error := errorInvoice{}
	json.NewDecoder(response.Body).Decode(&error)
	assert.Equal(t, []string{"sql: no rows in result set"}, error.Err, "should return error") //if not exists
}
