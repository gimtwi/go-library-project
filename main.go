package main

import (
	"github.com/gimtwi/go-library-project/controllers"
	"github.com/gimtwi/go-library-project/middleware"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gimtwi/go-library-project/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	utils.LoadENV()
	utils.ConnectDB()
	utils.MigrateDB()
}

func main() {
	r := gin.Default()

	bookRepo := types.NewBookRepository(utils.DB)
	userRepo := types.NewUserRepository(utils.DB)

	// user CRUD controller
	r.GET("/users", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.GetAllUsers(userRepo))
	r.GET("/users/:id", controllers.GetUserByID(userRepo))
	r.POST("/users/register", controllers.RegisterUser(userRepo))
	r.POST("/login", middleware.RateLimitMiddleware(), controllers.Login(userRepo))
	r.GET("/logout", controllers.Logout(userRepo))
	r.PUT("/users/new-moderator/:id", controllers.AssignRole(userRepo, types.Moderator))
	r.PUT("/users/new-admin/:id", controllers.AssignRole(userRepo, types.Admin))
	r.PUT("/users/:id/change-password", middleware.CheckIfMe(userRepo), controllers.ChangePassword(userRepo))
	r.DELETE("/users/:id", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.DeleteUser(userRepo))

	// book CRUD controller
	r.GET("/books", controllers.GetAllBooks(bookRepo))
	r.GET("/books/:id", controllers.GetBookByID(bookRepo))
	r.POST("/books", controllers.CreateBook(bookRepo))
	r.PUT("/books/:id", controllers.UpdateBook(bookRepo))
	r.DELETE("/books/:id", controllers.DeleteBook(bookRepo))

	r.Run()

}
