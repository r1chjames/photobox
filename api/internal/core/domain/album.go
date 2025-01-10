package domain

import (
	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID          string         `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        string         `json:"tags"`
	Metadata    datatypes.JSON `json:"metadata"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Thumnail    string         `json:"thumbnail"`
}
