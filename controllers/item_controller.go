package controllers

import (
	"net/http"
	"strconv"

	help "github.com/gimtwi/go-library-project/helpers"
	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredItemsByTitle(ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items, err := ir.GetAll(string(filters.Order), filters.Filter, filters.Limit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}

func GetItemByID(ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
			return
		}

		item, err := ir.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func GetItemsByAuthorID(ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
			return
		}

		items, err := ir.GetItemsByAuthor(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no items found"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func GetItemsByGenreID(ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
			return
		}

		items, err := ir.GetItemsByGenre(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no items found"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func GetItemsByKindID(ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kind id"})
			return
		}

		items, err := ir.GetItemsByKind(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no items found"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func CreateItem(ir types.ItemRepository, ar types.AuthorRepository, gr types.GenreRepository, kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item types.Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := help.CheckAuthorsGenresKinds(&item, ar, gr, kr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ir.Create(&item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	}
}

func UpdateItem(ir types.ItemRepository, ar types.AuthorRepository, gr types.GenreRepository, hr types.HoldRepository, lr types.LoanRepository, kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
			return
		}

		var item types.Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.ID = uint(id)

		i, err := ir.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}

		err = help.CheckAuthorsGenresKinds(&item, ar, gr, kr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = help.DisassociateAuthorsGenresKinds(&item, ir)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if item.Quantity != i.Quantity {
			if err := help.RearrangeHolds(item.ID, hr, lr, ir); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if err := ir.Update(&item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func DeleteItem(repo types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
