package types

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	BookID  uint      `json:"bookID" binding:"required"`
	UserID  uint      `json:"userID" binding:"required"`
	DueDate time.Time `json:"dueDate"`
}

type LoanRepository interface {
	Create(loan *Loan) error
	GetAll() ([]Loan, error)
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

func (l *LoanRepositoryImpl) GetAll() ([]Loan, error) {
	var loans []Loan
	if err := l.db.Find(&loans).Error; err != nil {
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
