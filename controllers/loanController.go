package controllers

import (
	"net/http"
	"strconv"
	"time"

	help "github.com/gimtwi/go-library-project/helpers"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetLoansByUserID(repo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		loans, err := repo.GetByUserID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch loans"})
			return
		}
		c.JSON(http.StatusOK, loans)
	}
}

func GetLoansByBookID(repo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		loans, err := repo.GetByBookID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch loans"})
			return
		}
		c.JSON(http.StatusOK, loans)
	}
}

func CreateLoan(loanRepo types.LoanRepository, bookRepo types.BookRepository, holdRepo types.HoldRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loan types.Loan
		if err := c.ShouldBindJSON(&loan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		book, err := bookRepo.GetByID(loan.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book not found"})
			return
		}

		userLoans, err := loanRepo.GetByUserID(loan.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(userLoans) >= 10 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user can't loan more than 10 books"})
			return
		}

		loans, err := loanRepo.GetByBookID(loan.BookID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if uint(len(loans)) >= book.Quantity {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book is not available"})
			return
		}

		holds, err := holdRepo.GetByBookID(loan.BookID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if uint(len(holds)) >= book.Quantity {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book is not available"})
			return
		}

		loan.CheckoutDate = time.Now()
		loan.ExpireDate = time.Now().Add(14 * 24 * time.Hour)  // * expires in 14 days
		loan.RenewableOn = time.Now().Add(11 * 24 * time.Hour) // * 3 days before loan expires

		if err := loanRepo.Create(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, loan)
	}
}

func ReturnTheBook(loanRepo types.LoanRepository, holdRepo types.HoldRepository, bookRepo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		loan, err := loanRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := loanRepo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := help.RearrangeHolds(loan.BookID, holdRepo, loanRepo, bookRepo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
