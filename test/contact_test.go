package test

import (
	"encoding/json"
	"fmt"
	"go-rbac-starter/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.CreateContactRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "08123456789",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.NotEmpty(t, responseBody.Data.ID)
}

func TestCreateContactUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := model.CreateContactRequest{
		FirstName: "John",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong-token")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestGetContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/contacts/%s", contact.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, contact.ID, responseBody.Data.ID)
	assert.Equal(t, contact.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, contact.LastName, responseBody.Data.LastName)
}

func TestGetContactNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/contacts/not-found-id", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdateContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	requestBody := model.UpdateContactRequest{
		FirstName: "Jane",
		LastName:  "Updated",
		Email:     "jane@example.com",
		Phone:     "08999999999",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/contacts/%s", contact.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Jane", responseBody.Data.FirstName)
	assert.Equal(t, "Updated", responseBody.Data.LastName)
}

func TestUpdateContactNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.UpdateContactRequest{
		FirstName: "Jane",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/contacts/not-found-id", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/contacts/%s", contact.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestDeleteContactNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodDelete, "/api/contacts/not-found-id", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestSearchContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	for i := 0; i < 5; i++ {
		contact := SeedContact(t, user.ID)
		// update with unique data
		db.Model(contact).Updates(map[string]interface{}{
			"first_name": "Contact" + strconv.Itoa(i),
			"email":      "contact" + strconv.Itoa(i) + "@example.com",
		})
	}

	request := httptest.NewRequest(http.MethodGet, "/api/contacts?page=1&size=10", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, int64(5), responseBody.Paging.TotalItem)
	assert.Equal(t, 1, responseBody.Paging.Page)
}

func TestSearchContactWithPagination(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	for i := 0; i < 20; i++ {
		contact := SeedContact(t, user.ID)
		db.Model(contact).Updates(map[string]interface{}{
			"first_name": "Contact" + strconv.Itoa(i),
			"email":      "contact" + strconv.Itoa(i) + "@example.com",
		})
	}

	request := httptest.NewRequest(http.MethodGet, "/api/contacts?page=2&size=5", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(4), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
}

func TestSearchContactWithFilter(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	for i := 0; i < 5; i++ {
		contact := SeedContact(t, user.ID)
		db.Model(contact).Updates(map[string]interface{}{
			"first_name": "Contact" + strconv.Itoa(i),
			"email":      "contact" + strconv.Itoa(i) + "@example.com",
		})
	}

	request := httptest.NewRequest(http.MethodGet, "/api/contacts?page=1&size=10&name=Contact1", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, int64(1), responseBody.Paging.TotalItem)
}
