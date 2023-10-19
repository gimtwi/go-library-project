package types

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Genre struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type GenreRepository interface {
	Create(genre *Genre) error
	GetAll() ([]Genre, error)
	GetByID(id uint) (*Genre, error)
	Update(genre *Genre) error
	Delete(id uint) error
	CheckGenre(id uint) error
}

type GenreRepositoryImpl struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &GenreRepositoryImpl{db}
}

func (g *GenreRepositoryImpl) Create(genre *Genre) error {
	return g.db.Create(genre).Error
}

func (g *GenreRepositoryImpl) GetAll() ([]Genre, error) {
	var genres []Genre
	if err := g.db.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

func (g *GenreRepositoryImpl) GetByID(id uint) (*Genre, error) {
	var genre Genre
	if err := g.db.First(&genre, id).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

func (g *GenreRepositoryImpl) Update(genre *Genre) error {
	return g.db.Save(genre).Error
}

func (g *GenreRepositoryImpl) Delete(id uint) error {
	return g.db.Delete(&Genre{}, id).Error
}

func (g *GenreRepositoryImpl) CheckGenre(id uint) error {
	_, err := g.GetByID(id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
