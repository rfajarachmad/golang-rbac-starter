package test

import (
	"encoding/json"
	"go-rbac-starter/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewerCannotCreateContact(t *testing.T) {
	ClearAll()
	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	requestBody := model.CreateContactRequest{
		FirstName: "John",
		LastName:  "Doe",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/contacts", strings.NewReader(string(bodyJSON)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestViewerCanReadContacts(t *testing.T) {
	ClearAll()
	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	request := httptest.NewRequest(http.MethodGet, "/api/contacts", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestViewerCannotUpdateContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	requestBody := model.UpdateContactRequest{
		FirstName: "Updated",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID, strings.NewReader(string(bodyJSON)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestViewerCannotDeleteContact(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")
	contact := SeedContact(t, user.ID)

	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	request := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID, nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestViewerCannotUpdateUser(t *testing.T) {
	ClearAll()
	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	requestBody := model.UpdateUserRequest{
		Name: "Updated Name",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJSON)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestViewerCanReadCurrentUser(t *testing.T) {
	ClearAll()
	viewer := SeedUserWithRole(t, "Viewer User", "viewer@example.com", "rahasia", "viewer")

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", viewer.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestUserCannotAccessAdminEndpoints(t *testing.T) {
	ClearAll()
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}

func TestAdminCanListAllUsers(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")
	SeedUser(t, "Regular User", "user@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/admin/users?page=1&size=10", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.UserResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, 2, len(responseBody.Data))
	assert.NotNil(t, responseBody.Paging)
	assert.Equal(t, int64(2), responseBody.Paging.TotalItem)
}

func TestAdminCanGetAnyUser(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodGet, "/api/admin/users/"+strconv.Itoa(user.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, "Regular User", responseBody.Data.Name)
	assert.Equal(t, "user", responseBody.Data.Role)
}

func TestAdminCanAssignRole(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")

	viewerRoleId := getRoleId("viewer")

	requestBody := model.AssignRoleRequest{
		RoleID: viewerRoleId,
	}
	bodyJSON, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPatch, "/api/admin/users/"+strconv.Itoa(user.ID)+"/role", strings.NewReader(string(bodyJSON)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, "viewer", responseBody.Data.Role)
}

func TestAdminCanDeleteAnyUser(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")
	user := SeedUser(t, "Regular User", "user@example.com", "rahasia")

	request := httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+strconv.Itoa(user.ID), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[bool])
	_ = json.Unmarshal(bytes, responseBody)

	assert.True(t, responseBody.Data)
}

func TestAdminCanListRoles(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")

	request := httptest.NewRequest(http.MethodGet, "/api/admin/roles", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.RoleResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, 3, len(responseBody.Data))
}

func TestAdminCanGetRole(t *testing.T) {
	ClearAll()
	admin := SeedUserWithRole(t, "Admin User", "admin@example.com", "rahasia", "admin")

	adminRoleId := getRoleId("admin")

	request := httptest.NewRequest(http.MethodGet, "/api/admin/roles/"+strconv.Itoa(adminRoleId), nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", admin.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.RoleResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, "admin", responseBody.Data.Name)
	assert.True(t, len(responseBody.Data.Permissions) > 0)
}

func TestRegisterDefaultsToUserRole(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:     "New User",
		Email:    "new@example.com",
		Password: "password123",
	}
	bodyJSON, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJSON)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, responseBody)

	assert.Equal(t, "user", responseBody.Data.Role)
}
