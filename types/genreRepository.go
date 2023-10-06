package types

import "time"

type Genre struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
