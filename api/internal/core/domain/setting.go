package domain

import "time"

type Setting struct {
	Key          string    `gorm:"primarykey" json:"key"`
	Value        string    `json:"value"`
	FriendlyName string    `json:"friendlyName"`
	Category     string    `json:"category"`
	Type         string    `json:"type"`
	Options      string    `json:"options"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Settings struct {
	Settings []Setting `binding:"required"`
}
