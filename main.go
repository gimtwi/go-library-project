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

	userRepo := types.NewUserRepository(utils.DB)
	bookRepo := types.NewBookRepository(utils.DB)
	authorRepo := types.NewAuthorRepository(utils.DB)
	genreRepo := types.NewGenreRepository(utils.DB)
	holdRepo := types.NewHoldRepository(utils.DB)
	loanRepo := types.NewLoanRepository(utils.DB)
	notificationRepo := types.NewNotificationRepository(utils.DB)

	// user CRUD controller
	r.GET("/user", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.GetAllUsers(userRepo))
	r.GET("/user/:id", controllers.GetUserByID(userRepo))
	r.PUT("/user/new-moderator/:id", controllers.AssignRole(userRepo, types.Moderator))
	r.PUT("/user/new-admin/:id", controllers.AssignRole(userRepo, types.Admin))
	r.PUT("/user/:id/change-password", middleware.CheckIfMe(userRepo), controllers.ChangePassword(userRepo))
	r.DELETE("/user/:id", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.DeleteUser(userRepo))

	r.POST("/register", controllers.RegisterUser(userRepo))
	r.POST("/login", middleware.RateLimitMiddleware(), controllers.Login(userRepo))
	r.GET("/logout", controllers.Logout(userRepo))

	// book CRUD controller
	r.GET("/book", controllers.GetAllBooks(bookRepo))
	r.GET("/book/:id", controllers.GetBookByID(bookRepo))
	r.POST("/book", controllers.CreateBook(bookRepo))
	r.PUT("/book/:id", controllers.UpdateBook(bookRepo))
	r.DELETE("/book/:id", controllers.DeleteBook(bookRepo))

	// author CRUD controller
	r.GET("/author", controllers.GetAllAuthors(authorRepo))
	r.GET("/author/:id", controllers.GetAuthorByID(authorRepo))
	r.POST("/author", controllers.CreateAuthor(authorRepo))
	r.PUT("/author/:id", controllers.UpdateAuthor(authorRepo))
	r.DELETE("/author/:id", controllers.DeleteAuthor(authorRepo))

	// genre CRUD controller
	r.GET("/genre", controllers.GetAllGenres(genreRepo))
	r.GET("/genre/:id", controllers.GetGenreByID(genreRepo))
	r.POST("/genre", controllers.CreateGenre(genreRepo))
	r.PUT("/genre/:id", controllers.UpdateGenre(genreRepo))
	r.DELETE("/genre/:id", controllers.DeleteGenre(genreRepo))

	// hold CRUD controller
	r.GET("/hold", controllers.GetAllHolds(holdRepo))
	r.GET("/hold/:id", controllers.GetHoldByID(holdRepo))
	r.POST("/hold", controllers.CreateHold(holdRepo))
	r.PUT("/hold/:id", controllers.UpdateHold(holdRepo))
	r.DELETE("/hold/:id", controllers.DeleteHold(holdRepo))

	// loan CRUD controller
	r.GET("/loan", controllers.GetAllLoans(loanRepo))
	r.GET("/loan/:id", controllers.GetLoanByID(loanRepo))
	r.POST("/loan", controllers.CreateLoan(loanRepo))
	r.PUT("/loan/:id", controllers.UpdateLoan(loanRepo))
	r.DELETE("/loan/:id", controllers.DeleteLoan(loanRepo))

	// notification CRUD controller
	r.GET("/notification", controllers.GetAllNotifications(notificationRepo))
	r.GET("/notification/:id", controllers.GetNotificationByID(notificationRepo))
	r.POST("/notification", controllers.CreateNotification(notificationRepo))
	r.PUT("/notification/:id", controllers.UpdateNotification(notificationRepo))
	r.DELETE("/notification/:id", controllers.DeleteNotification(notificationRepo))

	r.Run()

}
