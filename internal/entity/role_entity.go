package entity

import "time"

type Role struct {
	ID          int          `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string       `gorm:"column:name;uniqueIndex"`
	Description string       `gorm:"column:description"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (r *Role) TableName() string {
	return "roles"
}
