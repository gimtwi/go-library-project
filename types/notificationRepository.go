package types

import "time"

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	UserID uint `json:"userID" binding:"required"`
	Read   bool `json:"read"`
}
