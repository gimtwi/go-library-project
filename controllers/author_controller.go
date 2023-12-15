package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredAuthorsByName(ar types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authors, err := ar.GetAll(string(filters.Order), filters.Filter, filters.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, authors)
	}
}

func GetAuthorByID(ar types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		author, err := ar.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
			return
		}
		c.JSON(http.StatusOK, author)
	}
}

func CreateAuthor(ar types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var author types.Author
		if err := c.ShouldBindJSON(&author); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ar.Create(&author); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, author)
	}
}

func UpdateAuthor(ar types.AuthorRepository, ir types.ItemRepository) gin.HandlerFunc {
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

		_, err = ar.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
			return
		}

		if err := ar.Update(&author); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		associatedItems, err := ir.GetItemsByAuthor(author.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, item := range associatedItems {
			item.Authors = append(item.Authors, author)
			if err := ir.Update(&item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, author)
	}
}

func DeleteAuthor(ar types.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		if err := ar.Delete(uint(id)); err != nil {
			if strings.Contains(err.Error(), "foreign key constraint \"fk_item_authors_author\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete author with associated items"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Status(http.StatusNoContent)
	}
}
