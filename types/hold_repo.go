package types

import (
	"sort"
	"time"

	"gorm.io/gorm"
)

type Hold struct {
	ID     uint `gorm:"primarykey" json:"id"`
	ItemID uint `json:"itemID" binding:"required"`
	UserID uint `json:"userID"`

	PlacedDate time.Time `json:"placedDate"`

	IsAvailable          bool      `json:"isAvailable"`
	ExpiryDate           time.Time `json:"expiryDate"`
	InLinePosition       uint      `json:"inLinePosition"`       // * place in line
	EstimatedWeeksToWait uint      `json:"estimatedWeeksToWait"` // * approximate waiting days
	DeliveryDate         time.Time `json:"deliveryDate"`         // * deliver the hold after the date
}

type HoldRepository interface {
	Create(hold *Hold) error
	GetByUserID(userID uint) ([]Hold, error)
	GetByItemID(itemID uint) ([]Hold, error)
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

func (h *HoldRepositoryImpl) GetByUserID(userID uint) ([]Hold, error) {
	var holds []Hold
	if err := h.db.Find(&holds).Where("userID = ?", userID).Error; err != nil {
		return nil, err
	}

	// sort the holds by EstimatedWaitDays (lowest number comes first)
	sort.Slice(holds, func(i, j int) bool {
		return holds[i].EstimatedWeeksToWait < holds[j].EstimatedWeeksToWait
	})

	return holds, nil
}

func (h *HoldRepositoryImpl) GetByItemID(itemID uint) ([]Hold, error) {
	var holds []Hold
	if err := h.db.Find(&holds).Where("itemID = ?", itemID).Error; err != nil {
		return nil, err
	}

	// sort the holds by PlacedDate (earliest date first)
	sort.Slice(holds, func(i, j int) bool {
		return holds[i].PlacedDate.Before(holds[j].PlacedDate)
	})

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
