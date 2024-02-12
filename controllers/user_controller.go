package controllers

import (
	"net/http"
	"os"

	"github.com/gimtwi/go-library-project/middleware"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllUsers(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := ur.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		converted := make([]*types.UserResponse, len(users))

		for i, user := range users {
			converted[i] = user.ConvertToUserResponse()
		}

		c.JSON(http.StatusOK, converted)
	}
}

func GetUserByID(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := ur.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		res := user.ConvertToUserResponse()
		c.JSON(http.StatusOK, res)
	}
}

func RegisterUser(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user types.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := ur.GetByUniqueField("username", user.Username); err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "invalid username"})
			return
		}

		if _, err := ur.GetByUniqueField("email", user.Email); err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "invalid email"})
			return
		}

		if !user.IsValidEmail(user.Email) {
			c.JSON(http.StatusConflict, gin.H{"error": "invalid email"})
			return
		}

		hash, err := user.HashPassword(user.Password)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash the password"})
			return
		}

		user.ID = uuid.New().String()
		user.Role = types.Member
		user.Password = hash

		if ur.Create(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user is successfully registered!"})
	}
}

func Login(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := ur.GetByUniqueField("username", req.Username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
			return
		}

		errPWD := user.CheckPassword(req.Password)

		if errPWD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
			return
		}

		token, err := middleware.GenerateJWT(user)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(os.Getenv("COOKIE_NAME"), token, 3600*24*30, "", "", false, true) //secure false only on localhost, change to true in prod

		c.JSON(http.StatusOK, gin.H{"message": "user is successfully authenticated!"})
	}
}

func Logout(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(os.Getenv("COOKIE_NAME"), "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
	}
}

func AssignRole(ur types.UserRepository, role types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := ur.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		user.Role = role

		if err := ur.Update(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user was successfully assigned a new role!"})
	}
}

func ChangePassword(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ChangePasswordRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")

		user, err := ur.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		errPWD := user.CheckPassword(req.CurrentPassword)

		if errPWD != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid current password"})
			return
		}

		hash, err := user.HashPassword(req.NewPassword)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash the password"})
			return
		}

		user.Password = hash

		if err := ur.Update(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		Logout(ur)

		c.JSON(http.StatusOK, gin.H{"message": "password was changed successfully!"})

	}
}

func DeleteUser(ur types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := ur.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
