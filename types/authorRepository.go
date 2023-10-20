package types

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name string `json:"name" binding:"required"`
}

type AuthorRepository interface {
	Create(author *Author) error
	GetAll() ([]Author, error)
	GetByID(id uint) (*Author, error)
	Update(author *Author) error
	Delete(id uint) error
}

type AuthorRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &AuthorRepositoryImpl{db}
}

func (a *AuthorRepositoryImpl) Create(author *Author) error {
	return a.db.Create(author).Error
}

func (a *AuthorRepositoryImpl) GetAll() ([]Author, error) {
	var authors []Author
	if err := a.db.Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, nil
}

func (a *AuthorRepositoryImpl) GetByID(id uint) (*Author, error) {
	var author Author
	if err := a.db.First(&author, id).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

func (a *AuthorRepositoryImpl) Update(author *Author) error {
	return a.db.Save(author).Error
}

func (a *AuthorRepositoryImpl) Delete(id uint) error {
	return a.db.Delete(&Author{}, id).Error
}
