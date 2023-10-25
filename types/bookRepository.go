package types

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

type Book struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Authors     []Author `gorm:"many2many:book_authors" json:"authors" binding:"required"`
	Genres      []Genre  `gorm:"many2many:book_genres" json:"genres" binding:"required"`
	Quantity    uint     `json:"quantity"`
}

type BookRepository interface {
	Create(book *Book) error
	GetAll(order, filter string, limit uint) ([]Book, error)
	GetByID(id uint) (*Book, error)
	GetBooksByAuthor(authorID uint) ([]Book, error)
	GetBooksByGenre(genreID uint) ([]Book, error)
	Update(book *Book) error
	Delete(id uint) error
	DisassociateGenre(book *Book, genre *Genre) error
	DisassociateAuthor(book *Book, author *Author) error
}

type FilteredRequestBody struct {
	Order  Order  `json:"order"`
	Filter string `json:"filter"`
	Limit  uint   `json:"limit"`
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

func (br *BookRepositoryImpl) GetAll(order, filter string, limit uint) ([]Book, error) {
	var books []Book
	if err := br.db.Preload("Authors").Preload("Genres").Order("title "+order).Where("title LIKE ?", filter+"%").Limit(int(limit)).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (br *BookRepositoryImpl) GetByID(id uint) (*Book, error) {
	var book Book
	if err := br.db.Preload("Authors").Preload("Genres").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (br *BookRepositoryImpl) GetBooksByAuthor(authorID uint) ([]Book, error) {
	var books []Book
	if err := br.db.Joins("JOIN book_authors ON books.id = book_authors.book_id").
		Where("book_authors.author_id = ?", authorID).Preload("Genres").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (br *BookRepositoryImpl) GetBooksByGenre(genreID uint) ([]Book, error) {
	var books []Book
	if err := br.db.Joins("JOIN book_genres ON books.id = book_genres.book_id").
		Where("book_genres.genre_id = ?", genreID).Preload("Authors").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (br *BookRepositoryImpl) Update(book *Book) error {
	return br.db.Save(book).Error
}

func (br *BookRepositoryImpl) Delete(id uint) error {
	var book Book
	if err := br.db.Preload("Authors").Preload("Genres").First(&book, id).Error; err != nil {
		return err
	}

	for i := range book.Authors {
		br.db.Model(&book.Authors[i]).Association("Books").Delete(&book)
	}

	for i := range book.Genres {
		br.db.Model(&book.Genres[i]).Association("Books").Delete(&book)
	}

	if err := br.db.Delete(&book, id).Error; err != nil {
		return err
	}

	return nil
}

func (br *BookRepositoryImpl) DisassociateGenre(book *Book, genre *Genre) error {
	found := false
	for i, g := range book.Genres {
		if g.ID == genre.ID {
			found = true
			book.Genres = append(book.Genres[:i], book.Genres[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("genre is not associated with the book")
	}

	result := br.db.Model(book).Association("Genres").Delete(genre)

	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return fmt.Errorf("genre is not associated with the book")
		}
		return result
	}

	return nil
}

func (br *BookRepositoryImpl) DisassociateAuthor(book *Book, author *Author) error {
	found := false
	for i, a := range book.Authors {
		if a.ID == author.ID {
			found = true
			book.Authors = append(book.Authors[:i], book.Authors[i+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("author is not associated with the book")
	}

	result := br.db.Model(book).Association("Authors").Delete(author)

	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return fmt.Errorf("author is not associated with the book")
		}
		return result
	}

	return nil
}
