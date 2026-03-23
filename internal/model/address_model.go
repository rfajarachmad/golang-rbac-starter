package model

import "time"

type AddressResponse struct {
	ID         string    `json:"id,omitempty"`
	Street     string    `json:"street,omitempty"`
	City       string    `json:"city,omitempty"`
	Province   string    `json:"province,omitempty"`
	PostalCode string    `json:"postal_code,omitempty"`
	Country    string    `json:"country,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type CreateAddressRequest struct {
	UserId     int    `json:"-" validate:"required"`
	ContactId  string `json:"-" validate:"required,max=100"`
	Street     string `json:"street,omitempty" validate:"max=255"`
	City       string `json:"city,omitempty" validate:"max=100"`
	Province   string `json:"province,omitempty" validate:"max=100"`
	PostalCode string `json:"postal_code,omitempty" validate:"max=20"`
	Country    string `json:"country,omitempty" validate:"max=100"`
}

type UpdateAddressRequest struct {
	ID         string `json:"-" validate:"required,max=100"`
	UserId     int    `json:"-" validate:"required"`
	ContactId  string `json:"-" validate:"required,max=100"`
	Street     string `json:"street,omitempty" validate:"max=255"`
	City       string `json:"city,omitempty" validate:"max=100"`
	Province   string `json:"province,omitempty" validate:"max=100"`
	PostalCode string `json:"postal_code,omitempty" validate:"max=20"`
	Country    string `json:"country,omitempty" validate:"max=100"`
}

type GetAddressRequest struct {
	ID        string `json:"-" validate:"required,max=100"`
	UserId    int    `json:"-" validate:"required"`
	ContactId string `json:"-" validate:"required,max=100"`
}

type DeleteAddressRequest struct {
	ID        string `json:"-" validate:"required,max=100"`
	UserId    int    `json:"-" validate:"required"`
	ContactId string `json:"-" validate:"required,max=100"`
}

type ListAddressRequest struct {
	UserId    int    `json:"-" validate:"required"`
	ContactId string `json:"-" validate:"required,max=100"`
}
