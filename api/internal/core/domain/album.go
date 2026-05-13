package domain

import (
	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID           string         `gorm:"primarykey" json:"id"`
	Name         string         `json:"name" gorm:"uniqueIndex"`
	Description  string         `json:"description"`
	Tags         string         `json:"tags"`
	Metadata     datatypes.JSON `json:"metadata"`
	CreatedAt    time.Time      `json:"createdAt"`
	CreatedEpoch int64          `json:"createdEpoch" gorm:"index"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Thumbnail    string         `json:"thumbnail"`
	CoverPhotoId string         `json:"coverPhotoId"`
}
