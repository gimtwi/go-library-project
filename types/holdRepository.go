package types

import "time"

type Hold struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	BookID    uint `json:"bookID" binding:"required"`
	UserID    uint `json:"userID" binding:"required"`
	Fulfilled bool `json:"fulfilled"`
}
