package domain

import "time"

type UserRole string

const (
	ADMINISTRATOR UserRole = "administrator"
	VIEWER        UserRole = "viewer"
	CONTRIBUTOR   UserRole = "contributor"
)

type User struct {
	ID        string   `gorm:"primarykey" json:"id"`
	Username  string   `json:"username"`
	Role      UserRole `json:"role"`
	Password  string   `json:"password"`
	Email     string   `json:"email"`
	Approved  bool     `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
