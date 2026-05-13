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
	Username  string   `json:"username" gorm:"uniqueIndex"`
	Role      UserRole `json:"role"`
	Password  string   `json:"password"`
	Email     string   `json:"email" gorm:"uniqueIndex"`
	Approved  bool     `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
