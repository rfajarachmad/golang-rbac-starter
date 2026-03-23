package model

import "time"

type RoleResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
}

type PermissionResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type AssignRoleRequest struct {
	UserID int `json:"-" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}

type ListUsersRequest struct {
	Page int `json:"page" validate:"required,min=1"`
	Size int `json:"size" validate:"required,min=1,max=100"`
}

type GetAnyUserRequest struct {
	ID int `json:"-" validate:"required"`
}

type DeleteAnyUserRequest struct {
	ID int `json:"-" validate:"required"`
}
