package test

import (
	"encoding/json"
	"fmt"
	"go-rbac-starter/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Merdeka",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "10110",
		Country:    "Indonesia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/contacts/%s/addresses", contact.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.NotEmpty(t, responseBody.Data.ID)
}

func TestCreateAddressContactNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.CreateAddressRequest{
		Street: "Jalan Merdeka",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts/not-found-id/addresses", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestListAddresses(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)
	SeedAddress(t, contact.ID)
	SeedAddress(t, contact.ID)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/contacts/%s/addresses", contact.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 2, len(responseBody.Data))
}

func TestListAddressesContactNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/contacts/not-found-id/addresses", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestGetAddress(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)
	address := SeedAddress(t, contact.ID)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/contacts/%s/addresses/%s", contact.ID, address.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, address.ID, responseBody.Data.ID)
	assert.Equal(t, address.Street, responseBody.Data.Street)
	assert.Equal(t, address.City, responseBody.Data.City)
}

func TestGetAddressNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/contacts/%s/addresses/not-found-id", contact.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestUpdateAddress(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)
	address := SeedAddress(t, contact.ID)

	requestBody := model.UpdateAddressRequest{
		Street:     "Jalan Updated",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "40100",
		Country:    "Indonesia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/contacts/%s/addresses/%s", contact.ID, address.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Jalan Updated", responseBody.Data.Street)
	assert.Equal(t, "Bandung", responseBody.Data.City)
	assert.Equal(t, "Jawa Barat", responseBody.Data.Province)
}

func TestUpdateAddressNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	requestBody := model.UpdateAddressRequest{
		Street: "Jalan Updated",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/contacts/%s/addresses/not-found-id", contact.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteAddress(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)
	address := SeedAddress(t, contact.ID)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/contacts/%s/addresses/%s", contact.ID, address.ID), nil)
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

func TestDeleteAddressNotFound(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/contacts/%s/addresses/not-found-id", contact.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
