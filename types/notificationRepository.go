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
	GetAll() ([]Notification, error)
	GetByID(id uint) (*Notification, error)
	Update(notification *Notification) error
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

func (n *NotificationRepositoryImpl) GetAll() ([]Notification, error) {
	var notifications []Notification
	if err := n.db.Find(&notifications).Error; err != nil {
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

func (n *NotificationRepositoryImpl) Update(notification *Notification) error {
	return n.db.Save(notification).Error
}

func (n *NotificationRepositoryImpl) Delete(id uint) error {
	return n.db.Delete(&Notification{}, id).Error
}
