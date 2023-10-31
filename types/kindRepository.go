package types

import (
	"time"

	"gorm.io/gorm"
)

type Kind struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`

	Items []Item `gorm:"many2many:item_kinds" json:"items"`
}

type KindRepository interface {
	Create(kind *Kind) error
	GetAll(order, filter string, limit uint) ([]Kind, error)
	GetByID(id uint) (*Kind, error)
	Update(kind *Kind) error
	Delete(id uint) error
}

type kindRepositoryImpl struct {
	db *gorm.DB
}

func NewKindRepository(db *gorm.DB) KindRepository {
	return &kindRepositoryImpl{db}
}

func (k *kindRepositoryImpl) Create(kind *Kind) error {
	return k.db.Create(kind).Error
}

func (k *kindRepositoryImpl) GetAll(order, filter string, limit uint) ([]Kind, error) {
	var kinds []Kind
	if err := k.db.Preload("Items").Order("name "+order).Where("name LIKE ?", filter+"%").Limit(int(limit)).Find(&kinds).Error; err != nil {
		return nil, err
	}
	return kinds, nil
}

func (k *kindRepositoryImpl) GetByID(id uint) (*Kind, error) {
	var kind Kind
	if err := k.db.Preload("Items").First(&kind, id).Error; err != nil {
		return nil, err
	}
	return &kind, nil
}

func (k *kindRepositoryImpl) Update(kind *Kind) error {
	return k.db.Save(kind).Error
}

func (k *kindRepositoryImpl) Delete(id uint) error {
	return k.db.Delete(&Kind{}, id).Error
}
