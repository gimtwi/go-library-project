package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetOrderedFilteredKindsByName(kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters types.FilteredRequestBody
		if err := c.ShouldBindJSON(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		kinds, err := kr.GetAll(string(filters.Order), filters.Filter, filters.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, kinds)
	}
}

func GetKindByID(kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kind id"})
			return
		}

		kind, err := kr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "kind not found"})
			return
		}
		c.JSON(http.StatusOK, kind)
	}
}

func CreateKind(kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var kind types.Kind
		if err := c.ShouldBindJSON(&kind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := kr.Create(&kind); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, kind)
	}
}

func UpdateKind(kr types.KindRepository, ir types.ItemRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kind id"})
			return
		}

		var kind types.Kind
		if err := c.ShouldBindJSON(&kind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		kind.ID = uint(id)

		_, err = kr.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "kind not found"})
			return
		}

		if err := kr.Update(&kind); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		associatedItems, err := ir.GetItemsByKind(kind.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, item := range associatedItems {
			item.Kinds = append(item.Kinds, kind)
			if err := ir.Update(&item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, kind)
	}
}

func DeleteKind(kr types.KindRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kind id"})
			return
		}

		if err := kr.Delete(uint(id)); err != nil {
			if strings.Contains(err.Error(), "foreign key constraint \"fk_item_kinds_kind\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete kind with associated items"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Status(http.StatusNoContent)
	}
}
