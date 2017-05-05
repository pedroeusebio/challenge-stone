package controller

import (
	"app/shared/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	config := ParseJsonFile("../../../config/config.json")
	database.Connect(config.Database)
}

func TestUserPOST1(t *testing.T) {
	user := url.Values{}
	user.Set("name", "pedro")
	user.Add("password", "12")
	request, _ := http.NewRequest("POST", "/user", bytes.NewBufferString(user.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	UserPOST(response, request, nil)
	result := errorUser{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.NotEqual(t, "user_create", result.Err, "Ok response is expected")
}

func TestUserPOST2(t *testing.T) {
	user := url.Values{}
	user.Set("name", "pedro123410")
	user.Add("password", "algumasenha")
	request, _ := http.NewRequest("POST", "/user", bytes.NewBufferString(user.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	UserPOST(response, request, nil)
	result := successUser{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, "user_create", result.Success, "Ok response is expected")
}

func TestUserGET1(t *testing.T) {
	user := url.Values{}
	user.Set("page", "0")
	user.Add("length", "5")
	request, _ := http.NewRequest("GET", "/user?"+user.Encode(), nil)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	UserGET(response, request, nil)
	result := successUser{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, "user_getall", result.Success, "should return error")
	assert.True(t, 5 >= len(result.User), "should have the less or equal than 5")
}

//test wrong column
func TestUserGET2(t *testing.T) {
	user := url.Values{}
	user.Set("page", "0")
	user.Add("order", "[{\"Column\":\"asdasd\",\"Order\":\"desc\"}]")
	request, _ := http.NewRequest("GET", "/user?"+user.Encode(), nil)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	UserGET(response, request, nil)
	result := errorUser{}
	json.NewDecoder(response.Body).Decode(&result)
	assert.Equal(t, "pq: column \"asdasd\" does not exist", result.Err, "should return error")
}
