package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredGenresByName(gr types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		genres, err := gr.GetAll(string(filters.Order), filters.Filter, filters.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)
	}
}

func GetGenreByID(gr types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		genre, err := gr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		c.JSON(http.StatusOK, genre)
	}
}

func CreateGenre(gr types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var genre types.Genre
		if err := c.ShouldBindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := gr.Create(&genre); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, genre)
	}
}

func UpdateGenre(gr types.GenreRepository, ir types.ItemRepository) gin.HandlerFunc {
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

		_, err = gr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}

		if err := gr.Update(&genre); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		associatedItems, err := ir.GetItemsByGenre(genre.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, item := range associatedItems {
			item.Genres = append(item.Genres, genre)
			if err := ir.Update(&item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, genre)
	}
}

func DeleteGenre(gr types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		if err := gr.Delete(uint(id)); err != nil {
			if strings.Contains(err.Error(), "foreign key constraint \"fk_item_genres_genre\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete genre with associated items"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Status(http.StatusNoContent)
	}
}
