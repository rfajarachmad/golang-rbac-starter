package converter

import (
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
)

func ContactToResponse(contact *entity.Contact) *model.ContactResponse {
	return &model.ContactResponse{
		ID:        contact.ID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
