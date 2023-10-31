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
	itemRepo := types.NewItemRepository(utils.DB)
	authorRepo := types.NewAuthorRepository(utils.DB)
	genreRepo := types.NewGenreRepository(utils.DB)
	kindRepo := types.NewKindRepository(utils.DB)
	holdRepo := types.NewHoldRepository(utils.DB)
	loanRepo := types.NewLoanRepository(utils.DB)

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

	// item CRUD controller
	r.GET("/item", controllers.GetOrderedFilteredItemsByTitle(itemRepo))
	r.GET("/item/:id", controllers.GetItemByID(itemRepo))
	r.GET("/item/author/:id", controllers.GetItemsByAuthorID(itemRepo))
	r.GET("/item/genre/:id", controllers.GetItemsByGenreID(itemRepo))
	r.GET("/item/kind/:id", controllers.GetItemsByKindID(itemRepo))
	r.POST("/item", controllers.CreateItem(itemRepo, authorRepo, genreRepo, kindRepo))
	r.PUT("/item/:id", controllers.UpdateItem(itemRepo, authorRepo, genreRepo, holdRepo, loanRepo, kindRepo))
	r.DELETE("/item/:id", controllers.DeleteItem(itemRepo))

	// author CRUD controller
	r.GET("/author", controllers.GetOrderedFilteredAuthorsByName(authorRepo))
	r.GET("/author/:id", controllers.GetAuthorByID(authorRepo))
	r.POST("/author", controllers.CreateAuthor(authorRepo))
	r.PUT("/author/:id", controllers.UpdateAuthor(authorRepo, itemRepo))
	r.DELETE("/author/:id", controllers.DeleteAuthor(authorRepo))

	// genre CRUD controller
	r.GET("/genre", controllers.GetOrderedFilteredGenresByName(genreRepo))
	r.GET("/genre/:id", controllers.GetGenreByID(genreRepo))
	r.POST("/genre", controllers.CreateGenre(genreRepo))
	r.PUT("/genre/:id", controllers.UpdateGenre(genreRepo, itemRepo))
	r.DELETE("/genre/:id", controllers.DeleteGenre(genreRepo))

	// kind CRUD controller
	r.GET("/kind", controllers.GetOrderedFilteredKindsByName(kindRepo))
	r.GET("/kind/:id", controllers.GetKindByID(kindRepo))
	r.POST("/kind", controllers.CreateKind(kindRepo))
	r.PUT("/kind/:id", controllers.UpdateKind(kindRepo, itemRepo))
	r.DELETE("/kind/:id", controllers.DeleteKind(kindRepo))

	// hold CRUD controller
	r.GET("/hold/user/:id", controllers.GetHoldsByUserID(holdRepo, loanRepo, itemRepo))
	r.GET("/hold/item/:id", controllers.GetHoldsByItemID(holdRepo))
	r.POST("/hold", controllers.PlaceHold(holdRepo, loanRepo, itemRepo))
	r.DELETE("/cancel-hold/:id", controllers.CancelHold(holdRepo, loanRepo, itemRepo, userRepo))
	r.DELETE("/resolve-hold/:id", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.ResolveHold(holdRepo, loanRepo, itemRepo))

	// loan CRUD controller
	r.GET("/loan/item/:id", controllers.GetLoansByItemID(loanRepo))
	r.GET("/loan/user/:id", controllers.GetLoansByUserID(loanRepo))
	r.POST("/loan", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.CreateLoan(loanRepo, itemRepo, holdRepo))
	r.DELETE("/loan/:id", middleware.CheckPrivilege(userRepo, types.Moderator), controllers.ReturnTheItem(loanRepo, holdRepo, itemRepo))

	r.Run()

}
