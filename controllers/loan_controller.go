package controllers

import (
	"net/http"
	"strconv"
	"time"

	help "github.com/gimtwi/go-library-project/helpers"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetLoansByUserID(lr types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		loans, err := lr.GetByUserID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch loans"})
			return
		}
		c.JSON(http.StatusOK, loans)
	}
}

func GetLoansByItemID(lr types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		loans, err := lr.GetByItemID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch loans"})
			return
		}
		c.JSON(http.StatusOK, loans)
	}
}

func CreateLoan(lr types.LoanRepository, ir types.ItemRepository, hr types.HoldRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loan types.Loan
		if err := c.ShouldBindJSON(&loan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item, err := ir.GetByID(loan.ItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "item not found"})
			return
		}

		userLoans, err := lr.GetByUserID(loan.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(userLoans) >= 10 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user can't loan more than 10 items"})
			return
		}

		loans, err := lr.GetByItemID(loan.ItemID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if uint(len(loans)) >= item.Quantity {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "item is not available"})
			return
		}

		holds, err := hr.GetByItemID(loan.ItemID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if uint(len(holds)) >= item.Quantity {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "item is not available"})
			return
		}

		loan.CheckoutDate = time.Now()
		loan.ExpireDate = time.Now().Add(14 * 24 * time.Hour)  // * expires in 14 days
		loan.RenewableOn = time.Now().Add(11 * 24 * time.Hour) // * 3 days before loan expires

		if err := lr.Create(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, loan)
	}
}

func ProlongLoan(lr types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO enable for user to prolong the loan if there are no pending holds
	}
}

func ReturnTheItem(lr types.LoanRepository, hr types.HoldRepository, ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		loan, err := lr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := lr.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := help.RearrangeHolds(loan.ItemID, hr, lr, ir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
