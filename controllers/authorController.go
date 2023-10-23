package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredAuthorsByName(repo types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authors, err := repo.GetAll(string(filters.Order), filters.Filter, filters.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, authors)
	}
}

func GetAuthorByID(repo types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		author, err := repo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
			return
		}
		c.JSON(http.StatusOK, author)
	}
}

func CreateAuthor(repo types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var author types.Author
		if err := c.ShouldBindJSON(&author); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Create(&author); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, author)
	}
}

func UpdateAuthor(authorRepo types.AuthorRepository, bookRepo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		var author types.Author
		if err := c.ShouldBindJSON(&author); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		author.ID = uint(id)

		_, err = authorRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
			return
		}

		if err := authorRepo.Update(&author); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		associatedBooks, err := bookRepo.GetBooksByAuthor(author.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, book := range associatedBooks {
			book.Authors = append(book.Authors, author)
			if err := bookRepo.Update(&book); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, author)
	}
}

func DeleteAuthor(repo types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			if strings.Contains(err.Error(), "foreign key constraint \"fk_book_authors_author\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete author with associated books"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Status(http.StatusNoContent)
	}
}
