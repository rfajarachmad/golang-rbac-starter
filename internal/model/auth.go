package model

type Auth struct {
	ID          int      `json:"id"`
	RoleName    string   `json:"role"`
	Permissions []string `json:"permissions"`
}
