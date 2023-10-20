package controllers

import (
	"net/http"
	"strconv"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetAllBooks(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func GetBookByID(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		book, err := repo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func CreateBook(bookRepo types.BookRepository, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var book types.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		associatedAuthors := make([]types.Author, 0)
		associatedGenres := make([]types.Genre, 0)

		for _, author := range book.Author {
			a, err := authorRepo.GetByID(author.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "author not found"})
				return
			}
			associatedAuthors = append(associatedAuthors, *a)
		}

		for _, genre := range book.Genre {
			g, err := genreRepo.GetByID(genre.ID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "genre not found"})
				return
			}
			associatedGenres = append(associatedGenres, *g)
		}

		book.Author = associatedAuthors
		book.Genre = associatedGenres

		if book.Quantity >= 1 {
			book.IsAvailable = true
		} else {
			book.IsAvailable = false
		}

		if err := bookRepo.Create(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, book)
	}
}

func UpdateBook(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		var book types.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		book.ID = uint(id)

		if err := repo.Update(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func DeleteBook(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
