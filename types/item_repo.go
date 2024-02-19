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

type Item struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Authors     []Author `gorm:"many2many:item_authors" json:"authors"`
	Genres      []Genre  `gorm:"many2many:item_genres" json:"genres"`
	Kinds       []Kind   `gorm:"many2many:item_kinds" json:"kinds"`
	Quantity    uint     `json:"quantity"`
}

type ItemRepository interface {
	Create(i *Item) error
	GetAll(order, filter string, limit uint) ([]Item, error)
	GetByID(id uint) (*Item, error)
	GetItemsByAuthor(authorID uint) ([]Item, error)
	GetItemsByGenre(genreID uint) ([]Item, error)
	GetItemsByKind(kindID uint) ([]Item, error)
	Update(item *Item) error
	Delete(id uint) error
	DisassociateGenre(i *Item, genre *Genre) error
	DisassociateKind(i *Item, kind *Kind) error
	DisassociateAuthor(i *Item, author *Author) error
}

type FilteredRequestBody struct {
	Order  Order  `json:"order"`
	Filter string `json:"filter"`
	Limit  uint   `json:"limit"`
}

type ItemRepositoryImpl struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &ItemRepositoryImpl{db}
}

func (i *ItemRepositoryImpl) Create(item *Item) error {
	// Begin a transaction
	tx := i.db.Begin()

	// Create the item
	if err := tx.Create(item).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Associate Authors, Genres, and Kinds with the item
	if len(item.Authors) > 0 {
		if err := tx.Model(item).Association("Authors").Replace(item.Authors); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(item.Genres) > 0 {
		if err := tx.Model(item).Association("Genres").Replace(item.Genres); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(item.Kinds) > 0 {
		if err := tx.Model(item).Association("Kinds").Replace(item.Kinds); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (i *ItemRepositoryImpl) GetAll(order, filter string, limit uint) ([]Item, error) {
	var items []Item
	if err := i.db.Preload("Authors").Preload("Genres").Preload("Kinds").Order("title "+order).Where("title LIKE ?", filter+"%").Limit(int(limit)).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (i *ItemRepositoryImpl) GetByID(id uint) (*Item, error) {
	var item Item
	if err := i.db.Preload("Authors").Preload("Genres").Preload("Kinds").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (i *ItemRepositoryImpl) GetItemsByAuthor(authorID uint) ([]Item, error) {
	var items []Item
	if err := i.db.Joins("JOIN item_authors ON items.id = item_authors.item_id").
		Where("item_authors.author_id = ?", authorID).Preload("Genres").Preload("Kinds").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (i *ItemRepositoryImpl) GetItemsByGenre(genreID uint) ([]Item, error) {
	var items []Item
	if err := i.db.Joins("JOIN item_genres ON items.id = item_genres.item_id").
		Where("item_genres.genre_id = ?", genreID).Preload("Authors").Preload("Kinds").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (i *ItemRepositoryImpl) GetItemsByKind(kindID uint) ([]Item, error) {
	var items []Item
	if err := i.db.Joins("JOIN item_kinds ON items.id = item_kinds.item_id").
		Where("item_kinds.kind_id = ?", kindID).Preload("Authors").Preload("Genres").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (i *ItemRepositoryImpl) Update(item *Item) error {
	return i.db.Save(item).Error
}

func (i *ItemRepositoryImpl) Delete(id uint) error {
	var item Item
	if err := i.db.Preload("Authors").Preload("Genres").First(&item, id).Error; err != nil {
		return err
	}

	for in := range item.Authors {
		i.db.Model(&item.Authors[in]).Association("Items").Delete(&item)
	}

	for in := range item.Genres {
		i.db.Model(&item.Genres[in]).Association("Items").Delete(&item)
	}

	if err := i.db.Delete(&item, id).Error; err != nil {
		return err
	}

	return nil
}

func (i *ItemRepositoryImpl) DisassociateGenre(item *Item, genre *Genre) error {
	found := false
	for in, g := range item.Genres {
		if g.ID == genre.ID {
			found = true
			item.Genres = append(item.Genres[:in], item.Genres[in+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("genre is not associated with the item")
	}

	result := i.db.Model(item).Association("Genres").Delete(genre)

	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return fmt.Errorf("genre is not associated with the item")
		}
		return result
	}

	return nil
}

func (i *ItemRepositoryImpl) DisassociateKind(item *Item, kind *Kind) error {
	found := false
	for in, g := range item.Kinds {
		if g.ID == kind.ID {
			found = true
			item.Kinds = append(item.Kinds[:in], item.Kinds[in+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("kind is not associated with the item")
	}

	result := i.db.Model(item).Association("Kinds").Delete(kind)

	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kind is not associated with the item")
		}
		return result
	}

	return nil
}

func (i *ItemRepositoryImpl) DisassociateAuthor(item *Item, author *Author) error {
	found := false
	for in, a := range item.Authors {
		if a.ID == author.ID {
			found = true
			item.Authors = append(item.Authors[:in], item.Authors[in+1:]...)
			break
		}
	}

	if !found {
		return fmt.Errorf("author is not associated with the item")
	}

	result := i.db.Model(item).Association("Authors").Delete(author)

	if result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			return fmt.Errorf("author is not associated with the item")
		}
		return result
	}

	return nil
}
