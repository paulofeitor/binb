package model

import "time"

type Room struct {
	ID          int `gorm:"column:id"`
	Name        string
	PlaylistURL string
	CreatedAt   *time.Time
	Enabled     bool
}

func (r *Room) TableName() string {
	return "rooms"
}
