package types

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	ItemID uint   `json:"itemID" binding:"required"`
	UserID string `json:"userID" binding:"required"`

	CheckoutDate time.Time `json:"checkoutDate"` // * date of loan creation
	ExpireDate   time.Time `json:"expireDate"`   // * date of loan expiration
	RenewableOn  time.Time `json:"renewableOn"`  // ? 3 days before expiration date send notification
}

type LoanRepository interface {
	Create(loan *Loan) error
	GetByUserID(itemID string) ([]Loan, error)
	GetByItemID(itemID uint) ([]Loan, error)
	GetByID(id uint) (*Loan, error)
	Update(loan *Loan) error
	Delete(id uint) error
}

type LoanRepositoryImpl struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &LoanRepositoryImpl{db}
}

func (l *LoanRepositoryImpl) Create(loan *Loan) error {
	return l.db.Create(loan).Error
}

func (l *LoanRepositoryImpl) GetByUserID(userID string) ([]Loan, error) {
	var loans []Loan
	if err := l.db.Find(&loans).Where("userID = ?", userID).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (l *LoanRepositoryImpl) GetByItemID(itemID uint) ([]Loan, error) {
	var loans []Loan
	if err := l.db.Find(&loans).Where("itemID = ?", itemID).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (l *LoanRepositoryImpl) GetByID(id uint) (*Loan, error) {
	var loan Loan
	if err := l.db.First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (l *LoanRepositoryImpl) Update(loan *Loan) error {
	return l.db.Save(loan).Error
}

func (l *LoanRepositoryImpl) Delete(id uint) error {
	return l.db.Delete(&Loan{}, id).Error
}
