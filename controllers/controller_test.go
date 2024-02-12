package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
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

	db.AutoMigrate(&types.Author{}, &types.Genre{}, &types.Kind{}, &types.User{}, &types.Hold{}, &types.Loan{}, &types.Item{})

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

func TestControllers(t *testing.T) {
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

	userRepo := types.NewUserRepository(set)

	router := gin.Default()
	router.GET("/user", GetAllUsers(userRepo))
	router.GET("/user/:id", GetUserByID(userRepo))
	router.POST("/register", RegisterUser(userRepo))
	router.DELETE("/user/:id", DeleteUser(userRepo))

	t.Run("UserController", func(t *testing.T) {
		//* register user
		validUserJSON := []byte(`{
			"username": "test_user",
			"email": "user@test.com",
			"password": "password"
		}`)

		registerUser, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(validUserJSON))
		registerUser.Header.Set("Content-Type", "application/json")
		wRegisterUser := httptest.NewRecorder()
		router.ServeHTTP(wRegisterUser, registerUser)
		assert.Equal(t, http.StatusOK, wRegisterUser.Code)

		//* get all user
		getUsers, err := http.NewRequest("GET", "/user", nil)
		if err != nil {
			t.Fatal(err)
		}
		wUsers := httptest.NewRecorder()
		router.ServeHTTP(wUsers, getUsers)
		assert.Equal(t, http.StatusOK, wUsers.Code)

		var response []types.UserResponse
		err = json.Unmarshal(wUsers.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		userID := response[0].ID

		//* get user by id
		getUser, _ := http.NewRequest("GET", "/user/"+userID, nil)
		wUserID := httptest.NewRecorder()
		router.ServeHTTP(wUserID, getUser)
		assert.Equal(t, http.StatusOK, wUserID.Code)

		//* delete user
		deleteUser, _ := http.NewRequest("DELETE", "/user/"+userID, nil)
		wDeleteUser := httptest.NewRecorder()
		router.ServeHTTP(wDeleteUser, deleteUser)
		assert.Equal(t, http.StatusNoContent, wDeleteUser.Code)
	})

}
