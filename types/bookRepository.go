package types

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description"`
}

type BookRepository interface {
	Create(book *Book) error
	GetAll() ([]Book, error)
	GetByID(id uint) (*Book, error)
	Update(book *Book) error
	Delete(id uint) error
}

type BookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &BookRepositoryImpl{db}
}

func (br *BookRepositoryImpl) Create(book *Book) error {
	return br.db.Create(book).Error
}

func (br *BookRepositoryImpl) GetAll() ([]Book, error) {
	var books []Book
	if err := br.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (br *BookRepositoryImpl) GetByID(id uint) (*Book, error) {
	var book Book
	if err := br.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (br *BookRepositoryImpl) Update(book *Book) error {
	return br.db.Save(book).Error
}

func (br *BookRepositoryImpl) Delete(id uint) error {
	return br.db.Delete(&Book{}, id).Error
}
