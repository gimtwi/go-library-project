package controllers

import (
	"net/http"
	"strconv"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetNotificationByUserID(repo types.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
			return
		}

		notifications, err := repo.GetByUserID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		c.JSON(http.StatusOK, notifications)
	}
}

func CreateNotification(repo types.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification types.Notification
		if err := c.ShouldBindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Create(&notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, notification)
	}
}

func DeleteNotification(repo types.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
