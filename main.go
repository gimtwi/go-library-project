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
	// gin.SetMode(gin.ReleaseMode)

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
	r.PUT("/user/:id/change-password", middleware.CompareCookiesAndParameter(userRepo), controllers.ChangePassword(userRepo))
	r.DELETE("/user/:id", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.DeleteUser(userRepo))

	r.POST("/register", controllers.RegisterUser(userRepo))
	r.POST("/login", middleware.RateLimitMiddleware(), controllers.Login(userRepo))
	r.GET("/logout", controllers.Logout(userRepo))

	// TODO add middleware to each route

	// book CRUD controller
	r.GET("/book", controllers.GetOrderedFilteredBooksByTitle(bookRepo))
	r.GET("/book/:id", controllers.GetBookByID(bookRepo))
	r.GET("/book/author/:id", controllers.GetBooksByAuthorID(bookRepo))
	r.GET("/book/genre/:id", controllers.GetBooksByGenreID(bookRepo))
	r.POST("/book", controllers.CreateBook(bookRepo, authorRepo, genreRepo))
	r.PUT("/book/:id", controllers.UpdateBook(bookRepo, authorRepo, genreRepo))
	r.DELETE("/book/:id", controllers.DeleteBook(bookRepo))

	// author CRUD controller
	r.GET("/author", controllers.GetOrderedFilteredAuthorsByName(authorRepo))
	r.GET("/author/:id", controllers.GetAuthorByID(authorRepo))
	r.POST("/author", controllers.CreateAuthor(authorRepo))
	r.PUT("/author/:id", controllers.UpdateAuthor(authorRepo, bookRepo))
	r.DELETE("/author/:id", controllers.DeleteAuthor(authorRepo))

	// genre CRUD controller
	r.GET("/genre", controllers.GetOrderedFilteredGenresByName(genreRepo))
	r.GET("/genre/:id", controllers.GetGenreByID(genreRepo))
	r.POST("/genre", controllers.CreateGenre(genreRepo))
	r.PUT("/genre/:id", controllers.UpdateGenre(genreRepo, bookRepo))
	r.DELETE("/genre/:id", controllers.DeleteGenre(genreRepo))

	// hold CRUD controller
	r.GET("/hold/user/:id", controllers.GetHoldsByUserID(holdRepo, loanRepo, bookRepo))
	r.GET("/hold/book/:id", controllers.GetHoldsByBookID(holdRepo))
	r.POST("/hold", controllers.PlaceHold(holdRepo, loanRepo, bookRepo, userRepo))
	r.DELETE("/cancel-hold/:id", controllers.CancelHold(holdRepo, loanRepo, bookRepo))
	r.DELETE("/resolve-hold/:id", controllers.ResolveHold(holdRepo, loanRepo, bookRepo))

	// loan CRUD controller
	r.GET("/loan/book/:id", controllers.GetLoansByBookID(loanRepo))
	r.GET("/loan/user/:id", controllers.GetLoansByUserID(loanRepo))
	r.POST("/loan", controllers.CreateLoan(loanRepo, bookRepo))
	r.PUT("/loan/:id", controllers.UpdateLoan(loanRepo))
	r.DELETE("/loan/:id", controllers.ReturnTheBook(loanRepo))

	// notification CRUD controller
	r.GET("/notification", controllers.GetAllNotifications(notificationRepo))
	r.GET("/notification/:id", controllers.GetNotificationByID(notificationRepo))
	r.POST("/notification", controllers.CreateNotification(notificationRepo))
	r.PUT("/notification/:id", controllers.UpdateNotification(notificationRepo))
	r.DELETE("/notification/:id", controllers.DeleteNotification(notificationRepo))

	r.Run()

}
