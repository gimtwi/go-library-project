package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredGenresByName(repo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		genres, err := repo.GetAll(string(filters.Order), filters.Filter, filters.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)
	}
}

func GetGenreByID(repo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		genre, err := repo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		c.JSON(http.StatusOK, genre)
	}
}

func CreateGenre(repo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var genre types.Genre
		if err := c.ShouldBindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Create(&genre); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, genre)
	}
}

func UpdateGenre(genreRepo types.GenreRepository, bookRepo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		var genre types.Genre
		if err := c.ShouldBindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		genre.ID = uint(id)

		_, err = genreRepo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}

		if err := genreRepo.Update(&genre); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		associatedBooks, err := bookRepo.GetBooksByGenre(genre.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, book := range associatedBooks {
			book.Genres = append(book.Genres, genre)
			if err := bookRepo.Update(&book); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, genre)
	}
}

func DeleteGenre(repo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			if strings.Contains(err.Error(), "foreign key constraint \"fk_book_genres_genre\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete genre with associated books"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Status(http.StatusNoContent)
	}
}
