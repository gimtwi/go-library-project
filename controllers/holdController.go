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

func GetHoldsByUserID(hr types.HoldRepository, lr types.LoanRepository, ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		holds, err := hr.GetByUserID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch holds"})
			return
		}

		c.JSON(http.StatusOK, holds)
	}
}

func GetHoldsByItemID(hr types.HoldRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		holds, err := hr.GetByItemID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch holds"})
			return
		}
		c.JSON(http.StatusOK, holds)
	}
}

func PlaceHold(hr types.HoldRepository, lr types.LoanRepository, ir types.ItemRepository) gin.HandlerFunc {
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

		userHolds, err := hr.GetByUserID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		for _, h := range userHolds {
			if h.ItemID == hold.ItemID {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user has already placed a hold on this item"})
				return
			}
		}

		if len(userHolds) >= 10 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "users are not allowed to have more than 10 holds at the time"})
			return
		}

		hold.UserID = userID

		item, err := ir.GetByID(hold.ItemID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		holds, err := hr.GetByItemID(item.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch the holds for the item"})
			return
		}

		hold.PlacedDate = time.Now()
		hold.DeliveryDate = time.Now()

		if err := help.CheckInitialAvailability(&hold, item, hr, lr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if hold.IsAvailable {
			hold.InLinePosition = uint(len(holds) + 1)
			hold.ExpiryDate = time.Now().Add(3 * 24 * time.Hour) // if item is available for loner than 3 days the hold will expire automatically
			hold.EstimatedWeeksToWait = 0
		} else {
			hold.InLinePosition = uint(len(holds) + 1)
			if err := help.CalculateDaysToWait(&hold, item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if err := hr.Create(&hold); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, hold)
	}
}

// * by design supposed to be used when the hold is available to the user, but user wants to postpone the delivery
func ChangeDeliveryDate(hr types.HoldRepository) gin.HandlerFunc {
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

		updatedHold, err := hr.GetByID(hold.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "this hold doesn't exist"})
			return
		}

		updatedHold.DeliveryDate = hold.DeliveryDate

		if err := hr.Update(updatedHold); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//TODO update in line position
		//TODO update is available to false
		//TODO update other holds

		c.JSON(http.StatusOK, updatedHold)
	}

}

func CancelHold(hr types.HoldRepository, lr types.LoanRepository, ir types.ItemRepository, ur types.UserRepository) gin.HandlerFunc {
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

		hold, err := hr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold doesn't exist"})
			return
		}

		user, err := ur.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't find the user"})
			return
		}

		requiredRole := types.Moderator

		if hold.UserID == userID || types.CheckPrivilege(user.Role, requiredRole) {

			if err := hr.Delete(uint(id)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if err := help.RearrangeHolds(hold.ItemID, hr, lr, ir); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you can't perform this action"})
			return
		}

	}
}

// * must be performed by moderator
func ResolveHold(hr types.HoldRepository, lr types.LoanRepository, ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hold id"})
			return
		}

		hold, err := hr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold doesn't exist"})
			return
		}

		if !hold.IsAvailable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hold is not available"})
			return
		}

		var loan types.Loan

		loan.ItemID = hold.ItemID
		loan.UserID = hold.UserID
		loan.CheckoutDate = time.Now()
		loan.ExpireDate = time.Now().Add(14 * 24 * time.Hour)  // * expires in 14 days
		loan.RenewableOn = time.Now().Add(11 * 24 * time.Hour) // * 3 days before loan expires

		if err := lr.Create(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := hr.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := help.RearrangeHolds(hold.ItemID, hr, lr, ir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
