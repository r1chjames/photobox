package domain

import (
	"time"
)

type SharedLink struct {
	Token        string     `gorm:"primaryKey;size:64" json:"token"`
	ResourceType string     `gorm:"size:20;not null" json:"resourceType"`
	ResourceId   string     `gorm:"size:255;not null" json:"resourceId"`
	CreatedBy    string     `gorm:"size:36" json:"createdBy"`
	Expiry       *time.Time `json:"expiry"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	ViewCount    int        `gorm:"default:0" json:"viewCount"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"-"`
}
