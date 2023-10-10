package types

import (
	"time"

	"gorm.io/gorm"
)

type Hold struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	BookID    uint `json:"bookID" binding:"required"`
	UserID    uint `json:"userID" binding:"required"`
	Fulfilled bool `json:"fulfilled"`
}

type HoldRepository interface {
	Create(hold *Hold) error
	GetAll() ([]Hold, error)
	GetByID(id uint) (*Hold, error)
	Update(hold *Hold) error
	Delete(id uint) error
}

type HoldRepositoryImpl struct {
	db *gorm.DB
}

func NewHoldRepository(db *gorm.DB) HoldRepository {
	return &HoldRepositoryImpl{db}
}

func (h *HoldRepositoryImpl) Create(hold *Hold) error {
	return h.db.Create(hold).Error
}

func (h *HoldRepositoryImpl) GetAll() ([]Hold, error) {
	var holds []Hold
	if err := h.db.Find(&holds).Error; err != nil {
		return nil, err
	}
	return holds, nil
}

func (h *HoldRepositoryImpl) GetByID(id uint) (*Hold, error) {
	var hold Hold
	if err := h.db.First(&hold, id).Error; err != nil {
		return nil, err
	}
	return &hold, nil
}

func (h *HoldRepositoryImpl) Update(hold *Hold) error {
	return h.db.Save(hold).Error
}

func (h *HoldRepositoryImpl) Delete(id uint) error {
	return h.db.Delete(&Hold{}, id).Error
}
