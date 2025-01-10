package domain

import "time"

type Job struct {
	Name      string    `gorm:"primarykey" json:"name"`
	Status    string    `json:"status"`
	LastRun   time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
