package entity

import "time"

type Contact struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserId    int       `gorm:"column:user_id"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreignKey:user_id;references:id"`
	Addresses []Address `gorm:"foreignKey:contact_id;references:id"`
}

func (c *Contact) TableName() string {
	return "contacts"
}
