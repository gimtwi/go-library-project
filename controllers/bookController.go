package controllers

import (
	"net/http"
	"strconv"

	help "github.com/gimtwi/go-library-project/helpers"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredBooksByTitle(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		books, err := repo.GetAll(string(filters.Order), filters.Filter, filters.Limit)

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

func GetBooksByAuthorID(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		books, err := repo.GetBooksByAuthor(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no books found"})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func GetBooksByGenreID(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		books, err := repo.GetBooksByGenre(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no books found"})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func CreateBook(bookRepo types.BookRepository, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var book types.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := help.CheckAuthorsAndGenres(&book, authorRepo, genreRepo)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := bookRepo.Create(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, book)
	}
}

func UpdateBook(bookRepo types.BookRepository, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) gin.HandlerFunc {
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

		_, err = bookRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}

		err = help.CheckAuthorsAndGenres(&book, authorRepo, genreRepo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = help.DisassociateAuthorsAndGenres(&book, bookRepo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := bookRepo.Update(&book); err != nil {
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
