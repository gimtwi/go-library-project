package controllers

import (
	"net/http"
	"strconv"
	"time"

	help "github.com/gimtwi/go-library-project/helpers"
	"github.com/gimtwi/go-library-project/middleware"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetHoldsByUserID(repo types.HoldRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		holds, err := repo.GetByUserID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch holds"})
			return
		}
		c.JSON(http.StatusOK, holds)
	}
}

func GetHoldsByBookID(repo types.HoldRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		holds, err := repo.GetByBookID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch holds"})
			return
		}
		c.JSON(http.StatusOK, holds)
	}
}

func PlaceHold(holdRepo types.HoldRepository, loanRepo types.LoanRepository, bookRepo types.BookRepository, userRepo types.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var hold types.Hold
		if err := c.ShouldBindJSON(&hold); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := middleware.GetUserIDFromTheToken(c)
		if userID == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		hold.UserID = userID

		book, err := bookRepo.GetByID(hold.BookID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		holds, err := holdRepo.GetByBookID(book.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch the holds for the book"})
			return
		}

		loans, err := loanRepo.GetByBookID(book.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch the loans for the book"})
			return
		}

		hold.PlacedDate = time.Now()

		if len(holds) != 0 {

			if err := help.AssignAndUpdateHoldsCount(&hold, holds, holdRepo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			hold.HoldListPosition = hold.HoldsCount + 1

			holdsRatio := help.CalculateHoldsRatio(hold.OwnedCopies, hold.HoldsCount)

			if err := help.UpdateHoldsRatio(holds, holdsRatio, holdRepo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if err := help.CalculateEstimatedWaitDays(&hold); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		} else {
			hold.HoldsCount = 1
			hold.HoldListPosition = 1
			hold.HoldsRatio = 0
			hold.EstimatedWaitDays = 0
		}

		if len(loans) != 0 {

			if err := help.HandleOwnedCopies(&hold, loans); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		} else {
			hold.OwnedCopies = 0
		}

		if err := help.CalculateAvailableCopies(book, &hold, loans); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := holdRepo.Create(&hold); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, hold)
	}
}

func CancelHold(holdRepo types.HoldRepository, loanRepo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		userID := middleware.GetUserIDFromTheToken(c)
		if userID == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		hold, err := holdRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold doesn't exist"})
			return
		}

		if hold.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you can't perform this action"})
			return
		}

		if err := holdRepo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		holds, err := holdRepo.GetByBookID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch the holds for the book"})
			return
		}

		if len(holds) != 0 {

			if err := help.UpdateHoldsAfterDelete(uint(id), holds, holdRepo, loanRepo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}

	}
}

func ResolveHold(holdRepo types.HoldRepository, loanRepo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		hold, err := holdRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold doesn't exist"})
			return
		}

		if hold.HoldListPosition != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold is not available"})
			return
		}

		if !hold.IsAvailable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold is not available"})
			return
		}

		var loan types.Loan

		loan.BookID = hold.BookID
		loan.UserID = hold.UserID
		loan.CheckoutDate = time.Now()
		loan.ExpireDate = time.Now().Add(14 * 24 * time.Hour)  // * expires in 14 days
		loan.RenewableOn = time.Now().Add(11 * 24 * time.Hour) // * 3 days before loan expires

		if err := loanRepo.Create(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := holdRepo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		holds, err := holdRepo.GetByBookID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch the holds for the book"})
			return
		}

		if len(holds) != 0 {

			if err := help.UpdateHoldsAfterDelete(uint(id), holds, holdRepo, loanRepo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{"message": "hold has been resolved and new loan was created"})
	}
}
