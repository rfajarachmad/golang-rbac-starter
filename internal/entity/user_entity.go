package entity

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email;uniqueIndex"`
	Password  string    `gorm:"column:password"`
	Token     string    `gorm:"column:token"`
	RoleId    int       `gorm:"column:role_id"`
	Role      Role      `gorm:"foreignKey:role_id;references:id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Contacts  []Contact `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
