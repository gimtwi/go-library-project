package types

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_TEST_HOST"), os.Getenv("DB_TEST_PORT"), os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_PASSWORD"), os.Getenv("DB_TEST_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&Author{}, &Genre{}, &Kind{}, &User{}, &Hold{}, &Loan{}, &Item{})

	return db
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	set := setupTestDB()
	exitCode := m.Run()

	if sqlDB, err := set.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing test database: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "error getting underlying database connection: %v\n", err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func TestAuthorRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewAuthorRepository(set)

	author := &Author{Name: "TestAuthor"}

	t.Run("CreateAuthor", func(t *testing.T) {
		err := repo.Create(author)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, author.ID)

		defer func() {
			err := repo.Delete(author.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetAllAuthors", func(t *testing.T) {
		authors, err := repo.GetAll("asc", "", 1000)
		assert.NoError(t, err)
		assert.NotNil(t, authors)
	})

	t.Run("GetAuthorByID", func(t *testing.T) {
		err := repo.Create(author)
		assert.NoError(t, err)

		foundAuthor, err := repo.GetByID(author.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundAuthor)
		assert.Equal(t, author.Name, foundAuthor.Name)

		defer func() {
			err := repo.Delete(author.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateAuthor", func(t *testing.T) {
		err := repo.Create(author)
		assert.NoError(t, err)

		author.Name = "UpdatedAuthor"
		err = repo.Update(author)
		assert.NoError(t, err)

		updatedAuthor, err := repo.GetByID(author.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedAuthor", updatedAuthor.Name)

		defer func() {
			err := repo.Delete(author.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteAuthor", func(t *testing.T) {
		err := repo.Create(author)
		assert.NoError(t, err)

		err = repo.Delete(author.ID)
		assert.NoError(t, err)

		deletedAuthor, err := repo.GetByID(author.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedAuthor)
	})
}

func TestGenreRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewGenreRepository(set)

	genre := &Genre{Name: "TestGenre"}

	t.Run("CreateGenre", func(t *testing.T) {
		err := repo.Create(genre)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, genre.ID)

		defer func() {
			err := repo.Delete(genre.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetAllGenre", func(t *testing.T) {
		genre, err := repo.GetAll("asc", "", 1000)
		assert.NoError(t, err)
		assert.NotNil(t, genre)
	})

	t.Run("GetGenreByID", func(t *testing.T) {
		err := repo.Create(genre)
		assert.NoError(t, err)

		foundGenre, err := repo.GetByID(genre.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundGenre)
		assert.Equal(t, genre.Name, foundGenre.Name)

		defer func() {
			err := repo.Delete(genre.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateGenre", func(t *testing.T) {
		err := repo.Create(genre)
		assert.NoError(t, err)

		genre.Name = "UpdatedGenre"
		err = repo.Update(genre)
		assert.NoError(t, err)

		updatedGenre, err := repo.GetByID(genre.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedGenre", updatedGenre.Name)

		defer func() {
			err := repo.Delete(genre.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteGenre", func(t *testing.T) {
		err := repo.Create(genre)
		assert.NoError(t, err)

		err = repo.Delete(genre.ID)
		assert.NoError(t, err)

		deletedGenre, err := repo.GetByID(genre.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedGenre)
	})
}
func TestKindRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewKindRepository(set)

	kind := &Kind{Name: "TestKind"}

	t.Run("CreateKind", func(t *testing.T) {
		err := repo.Create(kind)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, kind.ID)

		defer func() {
			err := repo.Delete(kind.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetAllKind", func(t *testing.T) {
		kind, err := repo.GetAll("asc", "", 1000)
		assert.NoError(t, err)
		assert.NotNil(t, kind)
	})

	t.Run("GetKindByID", func(t *testing.T) {
		err := repo.Create(kind)
		assert.NoError(t, err)

		foundKind, err := repo.GetByID(kind.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundKind)
		assert.Equal(t, kind.Name, foundKind.Name)

		defer func() {
			err := repo.Delete(kind.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateKind", func(t *testing.T) {
		err := repo.Create(kind)
		assert.NoError(t, err)

		kind.Name = "UpdatedKind"
		err = repo.Update(kind)
		assert.NoError(t, err)

		updatedKind, err := repo.GetByID(kind.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedKind", updatedKind.Name)

		defer func() {
			err := repo.Delete(kind.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteKind", func(t *testing.T) {
		err := repo.Create(kind)
		assert.NoError(t, err)

		err = repo.Delete(kind.ID)
		assert.NoError(t, err)

		deletedKind, err := repo.GetByID(kind.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedKind)
	})
}

func TestUserRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewUserRepository(set)

	user := &User{
		FirstName: "TestFirstNameUser",
		LastName:  "TestLastNameUser",
		Username:  "test_username",
	}

	t.Run("CreateUser", func(t *testing.T) {
		err := repo.Create(user)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, user.ID)

		defer func() {
			err := repo.Delete(user.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetAllUser", func(t *testing.T) {
		user, err := repo.GetAll()
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		err := repo.Create(user)
		assert.NoError(t, err)

		foundUser, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.FirstName, foundUser.FirstName)
		assert.Equal(t, user.LastName, foundUser.LastName)

		defer func() {
			err := repo.Delete(user.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetUserByUniqueField", func(t *testing.T) {
		err := repo.Create(user)
		assert.NoError(t, err)

		foundUser, err := repo.GetByUniqueField("username", user.Username)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.Username, foundUser.Username)

		defer func() {
			err := repo.Delete(user.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateUser", func(t *testing.T) {
		err := repo.Create(user)
		assert.NoError(t, err)

		user.FirstName = "UpdatedFirstNameUser"
		err = repo.Update(user)
		assert.NoError(t, err)

		updatedUser, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedFirstNameUser", updatedUser.FirstName)

		defer func() {
			err := repo.Delete(user.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := repo.Create(user)
		assert.NoError(t, err)

		err = repo.Delete(user.ID)
		assert.NoError(t, err)

		deletedUser, err := repo.GetByID(user.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedUser)
	})
}

func TestHoldRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewHoldRepository(set)
	userRepo := NewUserRepository(set)
	itemRepo := NewItemRepository(set)

	user := &User{
		FirstName: "TestFirstNameUser",
		LastName:  "TestLastNameUser",
		Username:  "test_username",
	}

	item := &Item{Title: "TestTitle"}

	userErr := userRepo.Create(user)
	assert.NoError(t, userErr)
	assert.NotEqual(t, 0, user.ID)

	itemErr := itemRepo.Create(item)
	assert.NoError(t, itemErr)
	assert.NotEqual(t, 0, item.ID)

	hold := &Hold{
		ItemID: item.ID,
		UserID: user.ID,
	}

	t.Run("CreateHold", func(t *testing.T) {
		err := repo.Create(hold)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, hold.ID)

		defer func() {
			err := repo.Delete(hold.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetHoldsByUserID", func(t *testing.T) {
		hold, err := repo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, hold)
	})

	t.Run("GetHoldsByItemID", func(t *testing.T) {
		hold, err := repo.GetByItemID(item.ID)
		assert.NoError(t, err)
		assert.NotNil(t, hold)
	})

	t.Run("GetHoldByID", func(t *testing.T) {
		err := repo.Create(hold)
		assert.NoError(t, err)

		foundHold, err := repo.GetByID(hold.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundHold)
		assert.Equal(t, hold.UserID, foundHold.UserID)
		assert.Equal(t, hold.ItemID, foundHold.ItemID)

		defer func() {
			err := repo.Delete(hold.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateHold", func(t *testing.T) {
		err := repo.Create(hold)
		assert.NoError(t, err)

		time := time.Date(2023, time.December, 16, 9, 29, 7, 277199000, time.Local)
		hold.DeliveryDate = time
		err = repo.Update(hold)
		assert.NoError(t, err)

		updatedHold, err := repo.GetByID(hold.ID)
		assert.NoError(t, err)
		assert.Equal(t, time, updatedHold.DeliveryDate)

		defer func() {
			err := repo.Delete(hold.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteHold", func(t *testing.T) {
		err := repo.Create(hold)
		assert.NoError(t, err)

		err = repo.Delete(hold.ID)
		assert.NoError(t, err)

		deletedHold, err := repo.GetByID(hold.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedHold)
	})

	defer func() {
		userErr := userRepo.Delete(user.ID)
		itemErr := itemRepo.Delete(item.ID)
		assert.NoError(t, userErr)
		assert.NoError(t, itemErr)
	}()
}

func TestLoanRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewLoanRepository(set)
	userRepo := NewUserRepository(set)
	itemRepo := NewItemRepository(set)

	user := &User{
		FirstName: "TestFirstNameUser",
		LastName:  "TestLastNameUser",
		Username:  "test_username",
	}

	item := &Item{Title: "TestTitle"}

	userErr := userRepo.Create(user)
	assert.NoError(t, userErr)
	assert.NotEqual(t, 0, user.ID)

	itemErr := itemRepo.Create(item)
	assert.NoError(t, itemErr)
	assert.NotEqual(t, 0, item.ID)

	loan := &Loan{
		ItemID: item.ID,
		UserID: user.ID,
	}

	t.Run("CreateLoan", func(t *testing.T) {
		err := repo.Create(loan)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, loan.ID)

		defer func() {
			err := repo.Delete(loan.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetLoansByUserID", func(t *testing.T) {
		loan, err := repo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, loan)
	})

	t.Run("GetLoansByItemID", func(t *testing.T) {
		loan, err := repo.GetByItemID(item.ID)
		assert.NoError(t, err)
		assert.NotNil(t, loan)
	})

	t.Run("GetLoanByID", func(t *testing.T) {
		err := repo.Create(loan)
		assert.NoError(t, err)

		foundLoan, err := repo.GetByID(loan.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundLoan)
		assert.Equal(t, loan.UserID, foundLoan.UserID)
		assert.Equal(t, loan.ItemID, foundLoan.ItemID)

		defer func() {
			err := repo.Delete(loan.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateLoan", func(t *testing.T) {
		err := repo.Create(loan)
		assert.NoError(t, err)

		time := time.Date(2023, time.December, 16, 9, 29, 7, 277199000, time.Local)
		loan.ExpireDate = time
		err = repo.Update(loan)
		assert.NoError(t, err)

		updatedLoan, err := repo.GetByID(loan.ID)
		assert.NoError(t, err)
		assert.Equal(t, time, updatedLoan.ExpireDate)

		defer func() {
			err := repo.Delete(loan.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteLoan", func(t *testing.T) {
		err := repo.Create(loan)
		assert.NoError(t, err)

		err = repo.Delete(loan.ID)
		assert.NoError(t, err)

		deletedLoan, err := repo.GetByID(loan.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedLoan)
	})

	defer func() {
		userErr := userRepo.Delete(user.ID)
		itemErr := itemRepo.Delete(item.ID)
		assert.NoError(t, userErr)
		assert.NoError(t, itemErr)
	}()
}

func TestItemRepository(t *testing.T) {
	set := setupTestDB()

	defer func() {
		if sqlDB, err := set.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("error closing test database: %v", err)
			}
		} else {
			t.Errorf("error getting underlying database connection: %v", err)
		}
	}()

	repo := NewItemRepository(set)
	authorRepo := NewAuthorRepository(set)
	genreRepo := NewGenreRepository(set)
	kindRepo := NewKindRepository(set)

	author := &Author{Name: "TestAuthor"}
	genre := &Genre{Name: "TestGenre"}
	kind := &Kind{Name: "TestKind"}

	userErr := authorRepo.Create(author)
	assert.NoError(t, userErr)
	assert.NotEqual(t, 0, author.ID)

	genreErr := genreRepo.Create(genre)
	assert.NoError(t, genreErr)
	assert.NotEqual(t, 0, author.ID)

	kindErr := kindRepo.Create(kind)
	assert.NoError(t, kindErr)
	assert.NotEqual(t, 0, author.ID)

	item := &Item{
		Title:    "TestTitle",
		Authors:  []Author{{ID: author.ID}},
		Genres:   []Genre{{ID: genre.ID}},
		Kinds:    []Kind{{ID: kind.ID}},
		Quantity: 5,
	}

	t.Run("CreateItem", func(t *testing.T) {
		err := repo.Create(item)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, item.ID)

		defer func() {
			authErr := repo.DisassociateAuthor(item, author)
			assert.NoError(t, authErr)
			genreErr := repo.DisassociateGenre(item, genre)
			assert.NoError(t, genreErr)
			kindErr := repo.DisassociateKind(item, kind)
			assert.NoError(t, kindErr)
			err := repo.Delete(item.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("GetAllItems", func(t *testing.T) {
		item, err := repo.GetAll("asc", "", 1000)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("GetItemsByAuthorID", func(t *testing.T) {
		item, err := repo.GetItemsByAuthor(author.ID)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("GetItemsByGenreID", func(t *testing.T) {
		item, err := repo.GetItemsByGenre(genre.ID)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("GetItemsByKindID", func(t *testing.T) {
		item, err := repo.GetItemsByKind(kind.ID)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("GetItemByID", func(t *testing.T) {
		err := repo.Create(item)
		assert.NoError(t, err)

		foundItem, err := repo.GetByID(item.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundItem)
		assert.Equal(t, item.Quantity, foundItem.Quantity)

		defer func() {
			err := repo.Delete(item.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("UpdateItem", func(t *testing.T) {
		err := repo.Create(item)
		assert.NoError(t, err)

		item.Quantity = 6
		err = repo.Update(item)
		assert.NoError(t, err)

		updatedItem, err := repo.GetByID(item.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(6), updatedItem.Quantity)

		defer func() {
			err := repo.Delete(item.ID)
			assert.NoError(t, err)
		}()
	})

	t.Run("DeleteItem", func(t *testing.T) {
		err := repo.Create(item)
		assert.NoError(t, err)

		err = repo.Delete(item.ID)
		assert.NoError(t, err)

		deletedItem, err := repo.GetByID(item.ID)
		assert.Error(t, err)
		assert.Nil(t, deletedItem)
	})

	defer func() {
		userErr := authorRepo.Delete(author.ID)
		genreErr := genreRepo.Delete(genre.ID)
		kindErr := kindRepo.Delete(kind.ID)
		assert.NoError(t, userErr)
		assert.NoError(t, genreErr)
		assert.NoError(t, kindErr)
	}()
}
