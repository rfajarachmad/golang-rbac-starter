package entity

import "time"

type Permission struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;uniqueIndex"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (p *Permission) TableName() string {
	return "permissions"
}
