package controllers

import (
	"net/http"
	"strconv"

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

func CreateLoan(loanRepo types.LoanRepository, bookRepo types.BookRepository) gin.HandlerFunc {
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

		loans, err := loanRepo.GetByBookID(loan.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "loans not found"})
			return
		}

		if uint(len(loans)) >= book.Quantity {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "book is not available"})
			return
		}

		//TODO the rest of the logic

		if err := loanRepo.Create(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, loan)
	}
}

func UpdateLoan(repo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		var loan types.Loan
		if err := c.ShouldBindJSON(&loan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		loan.ID = uint(id)

		if err := repo.Update(&loan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, loan)
	}
}

func ReturnTheBook(repo types.LoanRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
