package types

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	UserID uint `json:"userID" binding:"required"`
	Read   bool `json:"read"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	GetByUserID(userID uint) ([]Notification, error)
	GetByID(id uint) (*Notification, error)
	Delete(id uint) error
}

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &NotificationRepositoryImpl{db}
}

func (n *NotificationRepositoryImpl) Create(notification *Notification) error {
	return n.db.Create(notification).Error
}

func (n *NotificationRepositoryImpl) GetByUserID(userID uint) ([]Notification, error) {
	var notifications []Notification
	if err := n.db.Find(&notifications).Where("userID = ?", userID).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (n *NotificationRepositoryImpl) GetByID(id uint) (*Notification, error) {
	var notification Notification
	if err := n.db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (n *NotificationRepositoryImpl) Delete(id uint) error {
	return n.db.Delete(&Notification{}, id).Error
}
