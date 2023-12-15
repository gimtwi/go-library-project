package types

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole uint

const (
	Admin     UserRole = 3
	Moderator UserRole = 2
	Member    UserRole = 1
)

type LoginRequest struct {
	Username string
	Password string
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	LibraryCard string   `gorm:"unique" json:"libraryCard"`
	Verified    string   `json:"verified"`
	Role        UserRole `json:"role"`

	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `gorm:"unique" json:"username"`
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	LibraryCard string   `json:"libraryCard"`
	Verified    string   `json:"verified"`
	Role        UserRole `json:"role"`

	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	PasswordExists bool   `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	Address        string `json:"address"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type UserRepository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	GetByUniqueField(field string, value string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (u *User) CheckPassword(password string) error {
	hashedPassword := []byte(u.Password)
	if len(hashedPassword) == 0 {
		return fmt.Errorf("hashed password is empty")
	}

	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func (u *User) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func (u *User) ConvertToUserResponse() *UserResponse {
	userResponse := UserResponse{
		ID:             u.ID,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Role:           u.Role,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Username:       u.Username,
		Email:          u.Email,
		PasswordExists: u.Password != "",
	}
	return &userResponse
}

func (ur *UserRepositoryImpl) Create(user *User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepositoryImpl) GetAll() ([]User, error) {
	var users []User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepositoryImpl) GetByID(id uint) (*User, error) {
	var user User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) GetByUniqueField(field string, value string) (*User, error) {
	var user User
	if err := ur.db.Where(field+" = ?", value).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) Update(user *User) error {
	return ur.db.Save(user).Error
}

func (ur *UserRepositoryImpl) Delete(id uint) error {
	return ur.db.Delete(&User{}, id).Error
}

func CheckPrivilege(userRole UserRole, privilegeRequired UserRole) bool {
	return userRole >= privilegeRequired
}
