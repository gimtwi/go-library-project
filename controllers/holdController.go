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

func GetHoldsByUserID(holdRepo types.HoldRepository, loanRepo types.LoanRepository, bookRepo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		holds, err := holdRepo.GetByUserID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch holds"})
			return
		}

		var holdsInfo []types.HoldsInfo

		for _, h := range holds {

			loans, err := loanRepo.GetByBookID(h.BookID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			book, err := bookRepo.GetByID(h.BookID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
				return
			}

			holdsCount := uint(len(holds))
			ownedCopies := uint(len(loans))
			holdsRatio := help.CalculateHoldsRatio(holdsCount, ownedCopies)

			holdsInfoItem := types.HoldsInfo{
				ID:                   h.ID,
				BookID:               h.BookID,
				UserID:               h.UserID,
				PlacedDate:           h.PlacedDate,
				IsAvailable:          h.IsAvailable,
				ExpiryDate:           h.ExpiryDate,
				InLinePosition:       h.InLinePosition,
				EstimatedWeeksToWait: h.EstimatedWeeksToWait,
				HoldsCount:           holdsCount,
				OwnedCopies:          ownedCopies,
				HoldsRatio:           holdsRatio,
				AvailableCopies:      book.Quantity - uint(len(loans)),
			}

			holdsInfo = append(holdsInfo, holdsInfoItem)
		}

		c.JSON(http.StatusOK, holdsInfo)
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

		hold.PlacedDate = time.Now()

		if err := help.CheckAvailability(&hold, book, holdRepo, loanRepo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if hold.IsAvailable {
			hold.InLinePosition = 0
			hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour) // if book is available for loner than 3 days the hold will expire automatically
			hold.EstimatedWeeksToWait = 0
		} else {
			hold.InLinePosition = uint(len(holds) + 1)
			if err := help.CalculateDaysToWait(&hold, book); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if err := holdRepo.Create(&hold); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, hold)
	}
}

func CancelHold(holdRepo types.HoldRepository, loanRepo types.LoanRepository, bookRepo types.BookRepository) gin.HandlerFunc {
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

		if err := help.RearrangeHolds(hold.BookID, holdRepo, loanRepo, bookRepo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
}

func ResolveHold(holdRepo types.HoldRepository, loanRepo types.LoanRepository, bookRepo types.BookRepository) gin.HandlerFunc {
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

		if err := help.RearrangeHolds(hold.BookID, holdRepo, loanRepo, bookRepo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
