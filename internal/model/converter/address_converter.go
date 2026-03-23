package converter

import (
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
)

func AddressToResponse(address *entity.Address) *model.AddressResponse {
	return &model.AddressResponse{
		ID:         address.ID,
		Street:     address.Street,
		City:       address.City,
		Province:   address.Province,
		PostalCode: address.PostalCode,
		Country:    address.Country,
		CreatedAt:  address.CreatedAt,
		UpdatedAt:  address.UpdatedAt,
	}
}
