package converter

import (
	"go-rbac-starter/internal/entity"
	"go-rbac-starter/internal/model"
)

func RoleToResponse(role *entity.Role) *model.RoleResponse {
	permissions := make([]model.PermissionResponse, len(role.Permissions))
	for i, p := range role.Permissions {
		permissions[i] = model.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		}
	}

	return &model.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
