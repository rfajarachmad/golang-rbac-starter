package test

import (
	"encoding/json"
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:     "Fajar Achmad",
		Email:    "fajar@example.com",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.NotZero(t, responseBody.Data.ID)
}

func TestRegisterValidationError(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:     "",
		Email:    "",
		Password: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.RegisterUserRequest{
		Name:     "Other User",
		Email:    "fajar@example.com",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestLogin(t *testing.T) {
	ClearAll()
	SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.LoginUserRequest{
		Email:    "fajar@example.com",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, responseBody.Data.Token)
}

func TestLoginWrongEmail(t *testing.T) {
	ClearAll()
	SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.LoginUserRequest{
		Email:    "wrong@example.com",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.LoginUserRequest{
		Email:    "fajar@example.com",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Name, responseBody.Data.Name)
	assert.Equal(t, user.Email, responseBody.Data.Email)
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong-token")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLogout(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
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

	// verify token is cleared
	updatedUser := new(entity.User)
	err = db.Where("id = ?", user.ID).First(updatedUser).Error
	assert.Nil(t, err)
	assert.Empty(t, updatedUser.Token)
}

func TestLogoutUnauthorized(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong-token")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.UpdateUserRequest{
		Name: "Fajar Updated",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Fajar Updated", responseBody.Data.Name)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Fajar Achmad", "fajar@example.com", "rahasia")

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	// verify password changed
	updatedUser := new(entity.User)
	err = db.Where("id = ?", user.ID).First(updatedUser).Error
	assert.Nil(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("rahasialagi"))
	assert.Nil(t, err)
}

func TestUpdateUnauthorized(t *testing.T) {
	ClearAll()

	requestBody := model.UpdateUserRequest{
		Name: "Updated",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong-token")

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
