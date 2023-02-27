package model

import "time"

type Timestampable struct {
	CreatedAt time.Time `gorm:"not null;autoCreateTime"    json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"    json:"updated_at"`
}
