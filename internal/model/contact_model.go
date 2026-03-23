package model

import "time"

type ContactResponse struct {
	ID        string    `json:"id,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreateContactRequest struct {
	UserId    int    `json:"-" validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"max=100"`
	Email     string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Phone     string `json:"phone,omitempty" validate:"max=100"`
}

type UpdateContactRequest struct {
	ID        string `json:"-" validate:"required,max=100"`
	UserId    int    `json:"-" validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"max=100"`
	Email     string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Phone     string `json:"phone,omitempty" validate:"max=100"`
}

type GetContactRequest struct {
	ID     string `json:"-" validate:"required,max=100"`
	UserId int    `json:"-" validate:"required"`
}

type DeleteContactRequest struct {
	ID     string `json:"-" validate:"required,max=100"`
	UserId int    `json:"-" validate:"required"`
}

type SearchContactRequest struct {
	UserId int    `json:"-" validate:"required"`
	Name   string `json:"name" validate:"max=100"`
	Email  string `json:"email" validate:"max=100"`
	Phone  string `json:"phone" validate:"max=100"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}
